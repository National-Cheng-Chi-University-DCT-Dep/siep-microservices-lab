package dto

import "github.com/google/uuid"

// 登入請求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johnsmith"`
	Password string `json:"password" binding:"required,min=6" example:"SecurePassword123!"`
}

// 註冊請求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"johnsmith"`
	Email    string `json:"email" binding:"required,email" example:"john.smith@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"SecurePassword123!"`
	FullName string `json:"full_name,omitempty" binding:"omitempty,max=100" example:"John Smith"`
}

// 刷新 Token 請求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
}

// 驗證信箱請求
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
}

// 更改密碼請求
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"OldPassword123!"`
	NewPassword     string `json:"new_password" binding:"required,min=8" example:"NewPassword456!"`
}

// 重設密碼請求
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// 完成重設密碼請求
type CompleteResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"NewPassword456!"`
}

// Supabase 登入回調請求
type SupabaseAuthCallback struct {
	SupabaseID   uuid.UUID `json:"supabase_id" binding:"required" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Email        string    `json:"email" binding:"required,email" example:"user@example.com"`
	Provider     string    `json:"provider" binding:"required" example:"google"`
	AccessToken  string    `json:"access_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	RefreshToken string    `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	UserMetadata struct {
		FullName  string `json:"full_name" example:"John Smith"`
		AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	} `json:"user_metadata"`
}

// 驗證 Token 請求
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
}

// 產生 API Key 請求
type GenerateAPIKeyRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=50" example:"My API Key"`
	Description string `json:"description" binding:"omitempty" example:"For integration with my app"`
	ExpiresIn   int    `json:"expires_in,omitempty" binding:"omitempty,min=1" example:"30"` // 天數
}

// 撤銷 API Key 請求
type RevokeAPIKeyRequest struct {
	KeyID string `json:"key_id" binding:"required,uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
}

// 確認重設密碼請求
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"NewPassword456!"`
}

// 更新個人資料請求
type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"omitempty,max=100" example:"John Smith"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url" example:"https://example.com/avatar.jpg"`
} 