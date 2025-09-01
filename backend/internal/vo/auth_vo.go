package vo

import (
	"time"

	"github.com/google/uuid"
)

// 認證回應
type AuthResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	ExpiresIn    int       `json:"expires_in" example:"3600"` // 秒數
	TokenType    string    `json:"token_type" example:"Bearer"`
	UserID       uuid.UUID `json:"user_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Username     string    `json:"username" example:"johnsmith"`
}

// 使用者資料回應
type UserProfileResponse struct {
	ID                    uuid.UUID  `json:"id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Username              string     `json:"username" example:"johnsmith"`
	Email                 string     `json:"email" example:"john.smith@example.com"`
	FullName              string     `json:"full_name,omitempty" example:"John Smith"`
	AvatarURL             string     `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	Role                  string     `json:"role" example:"premium"`
	IsActive              bool       `json:"is_active" example:"true"`
	EmailVerified         bool       `json:"email_verified" example:"true"`
	AuthProvider          string     `json:"auth_provider" example:"google"`
	SubscriptionType      string     `json:"subscription_type" example:"premium"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty" example:"2023-12-31T23:59:59Z"`
	APIQuota              int        `json:"api_quota" example:"1000"`
	APIUsage              int        `json:"api_usage" example:"50"`
	CreatedAt             time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	LastLogin             *time.Time `json:"last_login,omitempty" example:"2023-06-01T15:30:45Z"`
}

// API Key 回應
type APIKeyResponse struct {
	ID          uuid.UUID  `json:"id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Name        string     `json:"name" example:"My API Key"`
	Description string     `json:"description,omitempty" example:"For integration with my app"`
	Key         string     `json:"key,omitempty" example:"sk_live_************************"`
	Prefix      string     `json:"prefix" example:"sk_live_12345"`
	LastUsed    *time.Time `json:"last_used,omitempty" example:"2023-06-01T15:30:45Z"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" example:"2024-01-01T00:00:00Z"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	IsActive    bool       `json:"is_active" example:"true"`
}

// API Key 列表回應
type APIKeysListResponse struct {
	Keys []APIKeyResponse `json:"keys"`
}

// 驗證 Token 回應
type ValidateTokenResponse struct {
	Valid   bool      `json:"valid" example:"true"`
	UserID  uuid.UUID `json:"user_id,omitempty" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Expires time.Time `json:"expires,omitempty" example:"2023-07-01T00:00:00Z"`
}

// OAuth 提供者資訊
type OAuthProviderInfo struct {
	Provider    string `json:"provider" example:"google"`
	IsEnabled   bool   `json:"is_enabled" example:"true"`
	RedirectURI string `json:"redirect_uri,omitempty" example:"https://your-app.com/auth/callback/google"`
}

// OAuth 提供者列表回應
type OAuthProvidersResponse struct {
	Providers []OAuthProviderInfo `json:"providers"`
}

// 寄送驗證信結果回應
type EmailSentResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Verification email sent successfully"`
}

// 註冊回應
type RegisterResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User registered successfully"`
	UserID  uuid.UUID `json:"user_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
}

// 登入回應
type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login successful"`
	AuthResponse
}

// 刷新 Token 回應
type RefreshTokenResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Token refreshed successfully"`
	AuthResponse
}

// 登出回應
type LogoutResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Logout successful"`
}

// 更改密碼回應
type ChangePasswordResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Password changed successfully"`
}

// 重設密碼回應
type ResetPasswordResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Password reset email sent"`
}

// 確認重設密碼回應
type ConfirmResetPasswordResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Password reset confirmed"`
}

// 獲取使用者回應
type GetUserResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User retrieved successfully"`
	User    UserProfileResponse `json:"user"`
}

// 更新個人資料回應
type UpdateProfileResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Profile updated successfully"`
	User    UserProfileResponse `json:"user"`
}

// 認證 Token 回應
type AuthTokenResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	ExpiresIn    int       `json:"expires_in" example:"3600"`
	TokenType    string    `json:"token_type" example:"Bearer"`
	ExpiresAt    time.Time `json:"expires_at" example:"2023-07-01T00:00:00Z"`
	User         UserProfileResponse `json:"user"`
}

// 擴展使用者 VO
type ExtendedUserVO struct {
	UserProfileResponse
	APIKeys []APIKeyResponse `json:"api_keys,omitempty"`
} 