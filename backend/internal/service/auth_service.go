package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/model"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
	pkgjwt "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/jwt"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// AuthServiceInterface 認證服務介面
type AuthServiceInterface interface {
	Register(req *dto.RegisterRequest) (*vo.AuthTokenResponse, error)
	Login(req *dto.LoginRequest) (*vo.AuthTokenResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*vo.AuthTokenResponse, error)
	Logout(userID uuid.UUID, token string) error
	ChangePassword(userID uuid.UUID, req *dto.ChangePasswordRequest) error
	ResetPassword(req *dto.ResetPasswordRequest) error
	ConfirmResetPassword(req *dto.ConfirmResetPasswordRequest) error
	GetProfile(userID uuid.UUID) (*vo.ExtendedUserVO, error)
	UpdateProfile(userID uuid.UUID, req *dto.UpdateProfileRequest) (*vo.ExtendedUserVO, error)
	ValidateToken(token string) (*pkgjwt.JWTClaims, error)
}

// AuthService 認證服務實作
type AuthService struct {
	db         *gorm.DB
	jwtManager *pkgjwt.JWTManager
}

// NewAuthService 建立認證服務
func NewAuthService(db *gorm.DB, jwtManager *pkgjwt.JWTManager) AuthServiceInterface {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}

// Register 使用者註冊
func (s *AuthService) Register(req *dto.RegisterRequest) (*vo.AuthTokenResponse, error) {
	// 檢查使用者名稱是否已存在
	var existingUser model.User
	err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error
	if err == nil {
		if existingUser.Username == req.Username {
			return nil, dto.ErrUsernameExists
		}
		return nil, dto.ErrEmailExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		pkglogger.Error("Failed to check existing user", pkglogger.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		pkglogger.Error("Failed to hash password", pkglogger.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 設定預設角色
	role := "basic"

	// 建立使用者
	user := model.User{
		ID:                    uuid.New(),
		Username:              req.Username,
		Email:                 req.Email,
		PasswordHash:          string(hashedPassword),
		Role:                  model.UserRole(role),
		IsActive:              true,
		EmailVerified:         false,
		SubscriptionType:      role,
		SubscriptionExpiresAt: nil, // 基本用戶無期限
		APIQuota:              getDefaultAPIQuota(role),
		UsedAPIQuota:          0,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		pkglogger.Error("Failed to create user", pkglogger.Fields{
			"error":    err.Error(),
			"username": req.Username,
			"email":    req.Email,
		})
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 生成令牌
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Email, string(user.Role))
	if err != nil {
		pkglogger.Error("Failed to generate access token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		pkglogger.Error("Failed to generate refresh token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 更新最後登入時間
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// 轉換為 VO
	var userVO vo.UserVO
	if err := copier.Copy(&userVO, &user); err != nil {
		pkglogger.Error("Failed to copy user to VO", pkglogger.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to copy user data: %w", err)
	}

	response := &vo.AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1小時
		ExpiresAt:    time.Now().Add(time.Hour),
		User:         userVO,
	}

	pkglogger.Info("User registered successfully", pkglogger.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})

	return response, nil
}

// Login 使用者登入
func (s *AuthService) Login(req *dto.LoginRequest) (*vo.AuthTokenResponse, error) {
	// 驗證請求
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// 查找使用者
	var user model.User
	err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrInvalidCredentials
		}
		pkglogger.Error("Failed to find user", pkglogger.Fields{
			"error":    err.Error(),
			"username": req.Username,
		})
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// 檢查帳戶是否活躍
	if !user.IsActive {
		return nil, dto.ErrUserInactive
	}

	// 驗證密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, dto.ErrInvalidCredentials
	}

	// 生成令牌
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Email, string(user.Role))
	if err != nil {
		pkglogger.Error("Failed to generate access token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		pkglogger.Error("Failed to generate refresh token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 更新最後登入時間
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// 轉換為 VO
	var userVO vo.UserVO
	if err := copier.Copy(&userVO, &user); err != nil {
		pkglogger.Error("Failed to copy user to VO", pkglogger.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to copy user data: %w", err)
	}

	response := &vo.AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1小時
		ExpiresAt:    time.Now().Add(time.Hour),
		User:         userVO,
	}

	pkglogger.Info("User logged in successfully", pkglogger.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	})

	return response, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*vo.AuthTokenResponse, error) {
	// 驗證刷新令牌
	userID, err := s.jwtManager.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, dto.ErrInvalidRefreshToken
	}

	// 查找使用者
	var user model.User
	err = s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		pkglogger.Error("Failed to find user for refresh token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// 檢查帳戶是否活躍
	if !user.IsActive {
		return nil, dto.ErrUserInactive
	}

	// 生成新的令牌
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Email, string(user.Role))
	if err != nil {
		pkglogger.Error("Failed to generate new access token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		pkglogger.Error("Failed to generate new refresh token", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 轉換為 VO
	var userVO vo.UserVO
	if err := copier.Copy(&userVO, &user); err != nil {
		pkglogger.Error("Failed to copy user to VO", pkglogger.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to copy user data: %w", err)
	}

	response := &vo.AuthTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		ExpiresAt:    time.Now().Add(time.Hour),
		User:         userVO,
	}

	pkglogger.Info("Token refreshed successfully", pkglogger.Fields{
		"user_id": user.ID,
	})

	return response, nil
}

// Logout 使用者登出
func (s *AuthService) Logout(userID uuid.UUID, token string) error {
	// TODO: 實作令牌黑名單機制
	// 目前只記錄登出事件
	pkglogger.Info("User logged out", pkglogger.Fields{
		"user_id": userID,
	})
	return nil
}

// ChangePassword 修改密碼
func (s *AuthService) ChangePassword(userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	// 查找使用者
	var user model.User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// 驗證當前密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return dto.ErrInvalidCurrentPassword
	}

	// 加密新密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		pkglogger.Error("Failed to hash new password", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": userID,
		})
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密碼
	err = s.db.Model(&user).Updates(map[string]interface{}{
		"password_hash": string(hashedPassword),
		"updated_at":    time.Now(),
	}).Error
	if err != nil {
		pkglogger.Error("Failed to update password", pkglogger.Fields{
			"error":   err.Error(),
			"user_id": userID,
		})
		return fmt.Errorf("failed to update password: %w", err)
	}

	pkglogger.Info("Password changed successfully", pkglogger.Fields{
		"user_id": userID,
	})

	return nil
}

// ResetPassword 重設密碼
func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// TODO: 實作密碼重設邏輯
	// 1. 查找使用者
	// 2. 生成重設令牌
	// 3. 發送重設郵件
	pkglogger.Info("Password reset requested", pkglogger.Fields{
		"email": req.Email,
	})
	return nil
}

// ConfirmResetPassword 確認重設密碼
func (s *AuthService) ConfirmResetPassword(req *dto.ConfirmResetPasswordRequest) error {
	// TODO: 實作密碼重設確認邏輯
	// 1. 驗證重設令牌
	// 2. 更新密碼
	pkglogger.Info("Password reset confirmed", pkglogger.Fields{
		"token": req.Token[:8] + "...",
	})
	return nil
}

// GetProfile 取得個人檔案
func (s *AuthService) GetProfile(userID uuid.UUID) (*vo.ExtendedUserVO, error) {
	var user model.User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var userVO vo.ExtendedUserVO
	if err := copier.Copy(&userVO.UserVO, &user); err != nil {
		return nil, fmt.Errorf("failed to copy user data: %w", err)
	}

	// 設定額外欄位
	userVO.EmailVerified = user.EmailVerified
	userVO.SubscriptionType = user.SubscriptionType

	return &userVO, nil
}

// UpdateProfile 更新個人檔案
func (s *AuthService) UpdateProfile(userID uuid.UUID, req *dto.UpdateProfileRequest) (*vo.ExtendedUserVO, error) {
	var user model.User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// 更新欄位
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Username != nil {
		// 檢查使用者名稱是否已存在
		var existingUser model.User
		err := s.db.Where("username = ? AND id != ?", *req.Username, userID).First(&existingUser).Error
		if err == nil {
			return nil, dto.ErrUsernameExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check username uniqueness: %w", err)
		}
		updates["username"] = *req.Username
	}

	if req.Email != nil {
		// 檢查郵箱是否已存在
		var existingUser model.User
		err := s.db.Where("email = ? AND id != ?", *req.Email, userID).First(&existingUser).Error
		if err == nil {
			return nil, dto.ErrEmailExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		updates["email"] = *req.Email
		updates["email_verified"] = false // 需要重新驗證郵箱
	}

	// 執行更新
	err = s.db.Model(&user).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// 重新查詢更新後的使用者
	err = s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	var userVO vo.ExtendedUserVO
	if err := copier.Copy(&userVO.UserVO, &user); err != nil {
		return nil, fmt.Errorf("failed to copy user data: %w", err)
	}

	userVO.EmailVerified = user.EmailVerified
	userVO.SubscriptionType = user.SubscriptionType

	pkglogger.Info("Profile updated successfully", pkglogger.Fields{
		"user_id": userID,
	})

	return &userVO, nil
}

// ValidateToken 驗證令牌
func (s *AuthService) ValidateToken(token string) (*pkgjwt.JWTClaims, error) {
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	// 檢查使用者是否仍然存在且活躍
	var user model.User
	err = s.db.First(&user, "id = ?", claims.UserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, dto.ErrUserInactive
	}

	return claims, nil
}

// getDefaultAPIQuota 取得預設API配額
func getDefaultAPIQuota(role string) int {
	switch role {
	case "admin":
		return 100000 // 管理員無限制
	case "premium":
		return 10000 // 進階用戶10K
	default:
		return 1000 // 基本用戶1K
	}
} 