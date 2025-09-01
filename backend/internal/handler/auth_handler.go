package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
	pkgjwt "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/jwt"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	authService service.AuthServiceInterface
}

// NewAuthHandler 建立認證處理器
func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 使用者註冊
// @Summary 使用者註冊
// @Description 建立新的使用者帳戶
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "註冊資訊"
// @Success 201 {object} vo.RegisterResponse "註冊成功"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 409 {object} vo.BaseResponse "使用者名稱或郵箱已存在"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid register request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		handleServiceError(c, err, "Registration failed")
		return
	}

	response := vo.RegisterResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "User registered successfully",
			Timestamp: time.Now(),
		},
		Data: result,
	}

	c.JSON(http.StatusCreated, response)
}

// Login 使用者登入
// @Summary 使用者登入
// @Description 使用使用者名稱/郵箱和密碼進行登入
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登入資訊"
// @Success 200 {object} vo.LoginResponse "登入成功"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 401 {object} vo.BaseResponse "認證失敗"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid login request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		handleServiceError(c, err, "Login failed")
		return
	}

	response := vo.LoginResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Login successful",
			Timestamp: time.Now(),
		},
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken 刷新令牌
// @Summary 刷新存取令牌
// @Description 使用刷新令牌取得新的存取令牌
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} vo.RefreshTokenResponse "令牌刷新成功"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 401 {object} vo.BaseResponse "無效的刷新令牌"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid refresh token request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	result, err := h.authService.RefreshToken(&req)
	if err != nil {
		handleServiceError(c, err, "Token refresh failed")
		return
	}

	response := vo.RefreshTokenResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Token refreshed successfully",
			Timestamp: time.Now(),
		},
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

// Logout 使用者登出
// @Summary 使用者登出
// @Description 登出當前使用者並使令牌失效
// @Tags 認證
// @Security BearerAuth
// @Produce json
// @Success 200 {object} vo.LogoutResponse "登出成功"
// @Failure 401 {object} vo.BaseResponse "未授權"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	// 從標頭取得令牌
	authHeader := c.GetHeader("Authorization")
	token, err := pkgjwt.ExtractTokenFromHeader(authHeader)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_TOKEN", "Invalid authorization header", err)
		return
	}

	err = h.authService.Logout(userID.(uuid.UUID), token)
	if err != nil {
		handleServiceError(c, err, "Logout failed")
		return
	}

	response := vo.LogoutResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Logout successful",
			Timestamp: time.Now(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword 修改密碼
// @Summary 修改密碼
// @Description 修改當前使用者的密碼
// @Tags 認證
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "密碼修改資訊"
// @Success 200 {object} vo.ChangePasswordResponse "密碼修改成功"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 401 {object} vo.BaseResponse "未授權或當前密碼錯誤"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid change password request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	err := h.authService.ChangePassword(userID.(uuid.UUID), &req)
	if err != nil {
		handleServiceError(c, err, "Password change failed")
		return
	}

	response := vo.ChangePasswordResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Password changed successfully",
			Timestamp: time.Now(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword 重設密碼
// @Summary 重設密碼
// @Description 發送密碼重設郵件
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "重設密碼資訊"
// @Success 200 {object} vo.ResetPasswordResponse "重設郵件已發送"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid reset password request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	err := h.authService.ResetPassword(&req)
	if err != nil {
		handleServiceError(c, err, "Password reset failed")
		return
	}

	response := vo.ResetPasswordResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Password reset email sent",
			Timestamp: time.Now(),
		},
		Data: &struct {
			Message string `json:"message" example:"Password reset email sent"`
		}{
			Message: "If the email exists, you will receive password reset instructions",
		},
	}

	c.JSON(http.StatusOK, response)
}

// ConfirmResetPassword 確認重設密碼
// @Summary 確認重設密碼
// @Description 使用重設令牌確認新密碼
// @Tags 認證
// @Accept json
// @Produce json
// @Param request body dto.ConfirmResetPasswordRequest true "確認重設密碼資訊"
// @Success 200 {object} vo.ConfirmResetPasswordResponse "密碼重設成功"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 401 {object} vo.BaseResponse "無效的重設令牌"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/confirm-reset-password [post]
func (h *AuthHandler) ConfirmResetPassword(c *gin.Context) {
	var req dto.ConfirmResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid confirm reset password request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	err := h.authService.ConfirmResetPassword(&req)
	if err != nil {
		handleServiceError(c, err, "Password reset confirmation failed")
		return
	}

	response := vo.ConfirmResetPasswordResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Password reset successful",
			Timestamp: time.Now(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile 取得個人檔案
// @Summary 取得個人檔案
// @Description 取得當前使用者的個人檔案資訊
// @Tags 使用者
// @Security BearerAuth
// @Produce json
// @Success 200 {object} vo.GetUserResponse "取得成功"
// @Failure 401 {object} vo.BaseResponse "未授權"
// @Failure 404 {object} vo.BaseResponse "使用者不存在"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	result, err := h.authService.GetProfile(userID.(uuid.UUID))
	if err != nil {
		handleServiceError(c, err, "Failed to get profile")
		return
	}

	response := vo.GetUserResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Profile retrieved successfully",
			Timestamp: time.Now(),
		},
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile 更新個人檔案
// @Summary 更新個人檔案
// @Description 更新當前使用者的個人檔案資訊
// @Tags 使用者
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "更新資訊"
// @Success 200 {object} vo.UpdateProfileResponse "更新成功"
// @Failure 400 {object} vo.BaseResponse "請求參數錯誤"
// @Failure 401 {object} vo.BaseResponse "未授權"
// @Failure 409 {object} vo.BaseResponse "使用者名稱或郵箱已存在"
// @Failure 500 {object} vo.BaseResponse "內部伺服器錯誤"
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkglogger.Error("Invalid update profile request", pkglogger.Fields{
			"error": err.Error(),
		})
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request parameters", err)
		return
	}

	result, err := h.authService.UpdateProfile(userID.(uuid.UUID), &req)
	if err != nil {
		handleServiceError(c, err, "Profile update failed")
		return
	}

	response := vo.UpdateProfileResponse{
		BaseResponse: vo.BaseResponse{
			Success:   true,
			Message:   "Profile updated successfully",
			Timestamp: time.Now(),
		},
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes 註冊認證路由
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		// 公開路由
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/confirm-reset-password", h.ConfirmResetPassword)

		// 需要認證的路由
		// 這些路由需要在主程式中添加JWT中介軟體
	}

	// 個人檔案路由（需要認證）
	profile := router.Group("/auth")
	{
		profile.POST("/logout", h.Logout)
		profile.POST("/change-password", h.ChangePassword)
		profile.GET("/profile", h.GetProfile)
		profile.PUT("/profile", h.UpdateProfile)
	}
}

// respondError 回應錯誤
func respondError(c *gin.Context, statusCode int, code string, message string, err error) {
	errorVO := &vo.ErrorVO{
		Code:    code,
		Message: message,
	}
	
	if err != nil {
		errorVO.Details = err.Error()
	}

	response := vo.BaseResponse{
		Success:   false,
		Message:   message,
		Error:     errorVO,
		Timestamp: time.Now(),
	}
	c.JSON(statusCode, response)
}

// handleServiceError 處理服務錯誤
func handleServiceError(c *gin.Context, err error, message string) {
	pkglogger.Error("Service error", pkglogger.Fields{
		"error":   err.Error(),
		"message": message,
	})

	// 根據錯誤類型返回適當的狀態碼
	switch {
	case errors.Is(err, dto.ErrInvalidCredentials):
		respondError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password", err)
	case errors.Is(err, dto.ErrUserNotFound):
		respondError(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found", err)
	case errors.Is(err, dto.ErrUserInactive):
		respondError(c, http.StatusForbidden, "USER_INACTIVE", "User account is inactive", err)
	case errors.Is(err, dto.ErrUsernameExists):
		respondError(c, http.StatusConflict, "USERNAME_EXISTS", "Username already exists", err)
	case errors.Is(err, dto.ErrEmailExists):
		respondError(c, http.StatusConflict, "EMAIL_EXISTS", "Email already exists", err)
	case errors.Is(err, dto.ErrInvalidRefreshToken):
		respondError(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token", err)
	case errors.Is(err, dto.ErrInvalidCurrentPassword):
		respondError(c, http.StatusBadRequest, "INVALID_CURRENT_PASSWORD", "Invalid current password", err)
	default:
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, err)
	}
} 