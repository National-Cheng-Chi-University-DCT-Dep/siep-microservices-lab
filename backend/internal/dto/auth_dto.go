package dto

import "github.com/google/uuid"

// LoginRequest 登入請求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" validate:"required"`
	Password string `json:"password" binding:"required,min=6,max=128" validate:"required"`
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" validate:"required"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=128" validate:"required"`
	Role     string `json:"role" binding:"omitempty,oneof=basic premium" validate:"omitempty,oneof=basic premium"`
}

// RefreshTokenRequest 刷新令牌請求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}

// ChangePasswordRequest 修改密碼請求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=128" validate:"required"`
}

// ResetPasswordRequest 重設密碼請求
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

// ConfirmResetPasswordRequest 確認重設密碼請求
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128" validate:"required"`
}

// UpdateProfileRequest 更新個人檔案請求
type UpdateProfileRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email" validate:"omitempty,email"`
}

// APIKeyCreateRequest 建立API金鑰請求
type APIKeyCreateRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=100" validate:"required"`
	ExpiresAt *int64 `json:"expires_at" validate:"omitempty"`
	Quota     *int   `json:"quota" validate:"omitempty,min=1"`
}

// APIKeyUpdateRequest 更新API金鑰請求
type APIKeyUpdateRequest struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=100"`
	IsActive  *bool   `json:"is_active" validate:"omitempty"`
	ExpiresAt *int64  `json:"expires_at" validate:"omitempty"`
	Quota     *int    `json:"quota" validate:"omitempty,min=1"`
}

// APIKeyListRequest API金鑰列表請求
type APIKeyListRequest struct {
	Page     int  `json:"page" form:"page" validate:"omitempty,min=1"`
	PageSize int  `json:"page_size" form:"page_size" validate:"omitempty,min=1,max=100"`
	IsActive *bool `json:"is_active" form:"is_active" validate:"omitempty"`
}

// UserListRequest 使用者列表請求
type UserListRequest struct {
	Page     int     `json:"page" form:"page" validate:"omitempty,min=1"`
	PageSize int     `json:"page_size" form:"page_size" validate:"omitempty,min=1,max=100"`
	Role     *string `json:"role" form:"role" validate:"omitempty,oneof=admin premium basic"`
	IsActive *bool   `json:"is_active" form:"is_active" validate:"omitempty"`
	Search   *string `json:"search" form:"search" validate:"omitempty,max=100"`
}

// UserUpdateRequest 更新使用者請求（管理員用）
type UserUpdateRequest struct {
	Username              *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email                 *string `json:"email" validate:"omitempty,email"`
	Role                  *string `json:"role" validate:"omitempty,oneof=admin premium basic"`
	IsActive              *bool   `json:"is_active" validate:"omitempty"`
	SubscriptionExpiresAt *int64  `json:"subscription_expires_at" validate:"omitempty"`
	APIQuota              *int    `json:"api_quota" validate:"omitempty,min=0"`
}

// GetUserRequest 取得使用者請求
type GetUserRequest struct {
	UserID uuid.UUID `json:"user_id" uri:"user_id" binding:"required" validate:"required,uuid"`
}

// DeleteUserRequest 刪除使用者請求
type DeleteUserRequest struct {
	UserID uuid.UUID `json:"user_id" uri:"user_id" binding:"required" validate:"required,uuid"`
}

// ValidateCredentials 驗證登入憑證
func (r *LoginRequest) ValidateCredentials() error {
	if len(r.Username) < 3 {
		return ErrInvalidUsername
	}
	if len(r.Password) < 6 {
		return ErrInvalidPassword
	}
	return nil
}

// SetDefaults 設定預設值
func (r *APIKeyListRequest) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
}

// SetDefaults 設定預設值
func (r *UserListRequest) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
}

// GetOffset 取得分頁偏移量
func (r *APIKeyListRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// GetLimit 取得限制數量
func (r *APIKeyListRequest) GetLimit() int {
	return r.PageSize
}

// GetOffset 取得分頁偏移量
func (r *UserListRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// GetLimit 取得限制數量
func (r *UserListRequest) GetLimit() int {
	return r.PageSize
} 