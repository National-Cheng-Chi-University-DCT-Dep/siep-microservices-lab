package vo

import (
	"strings"
	"time"
)

// AuthTokenResponse 認證令牌回應
type AuthTokenResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string    `json:"token_type" example:"Bearer"`
	ExpiresIn    int       `json:"expires_in" example:"3600"`
	ExpiresAt    time.Time `json:"expires_at" example:"2024-12-01T15:00:00Z"`
	User         UserVO    `json:"user"`
}

// LoginResponse 登入回應
// @Description 成功登入後的回應資料
type LoginResponse struct {
	BaseResponse
	Data *AuthTokenResponse `json:"data,omitempty"`
}

// RegisterResponse 註冊回應
// @Description 成功註冊後的回應資料
type RegisterResponse struct {
	BaseResponse
	Data *AuthTokenResponse `json:"data,omitempty"`
}

// RefreshTokenResponse 刷新令牌回應
// @Description 刷新令牌後的回應資料
type RefreshTokenResponse struct {
	BaseResponse
	Data *AuthTokenResponse `json:"data,omitempty"`
}

// ExtendedUserVO 擴展的使用者資訊（包含認證相關額外欄位）
type ExtendedUserVO struct {
	UserVO
	EmailVerified    bool   `json:"email_verified" example:"true"`
	SubscriptionType string `json:"subscription_type" example:"premium"`
}

// GetUserResponse 取得使用者回應
// @Description 取得使用者資訊的回應
type GetUserResponse struct {
	BaseResponse
	Data *ExtendedUserVO `json:"data,omitempty"`
}

// UpdateProfileResponse 更新個人檔案回應
// @Description 更新個人檔案後的回應
type UpdateProfileResponse struct {
	BaseResponse
	Data *ExtendedUserVO `json:"data,omitempty"`
}

// UserListVO 使用者列表
type UserListVO struct {
	Users      []UserVO      `json:"users"`
	Pagination PaginationVO  `json:"pagination"`
}

// GetUserListResponse 取得使用者列表回應
// @Description 取得使用者列表的回應
type GetUserListResponse struct {
	BaseResponse
	Data *UserListVO `json:"data,omitempty"`
}

// ExtendedAPIKeyVO 擴展的API金鑰資訊（包含額外欄位）
type ExtendedAPIKeyVO struct {
	APIKeyVO
	Key         string  `json:"key,omitempty" example:"sk_live_..."`
	KeyPreview  string  `json:"key_preview" example:"sk_live_****_****_1234"`
	UsedQuota   int     `json:"used_quota" example:"1250"`
	LastUsedIP  *string `json:"last_used_ip,omitempty" example:"192.168.1.100"`
}

// CreateAPIKeyResponse 建立API金鑰回應
// @Description 建立API金鑰後的回應，只有在建立時會回傳完整的金鑰
type CreateAPIKeyResponse struct {
	BaseResponse
	Data *ExtendedAPIKeyVO `json:"data,omitempty"`
}

// GetAPIKeyResponse 取得API金鑰回應
// @Description 取得API金鑰資訊的回應
type GetAPIKeyResponse struct {
	BaseResponse
	Data *ExtendedAPIKeyVO `json:"data,omitempty"`
}

// APIKeyListVO API金鑰列表
type APIKeyListVO struct {
	APIKeys    []ExtendedAPIKeyVO `json:"api_keys"`
	Pagination PaginationVO       `json:"pagination"`
}

// GetAPIKeyListResponse 取得API金鑰列表回應
// @Description 取得API金鑰列表的回應
type GetAPIKeyListResponse struct {
	BaseResponse
	Data *APIKeyListVO `json:"data,omitempty"`
}

// UpdateAPIKeyResponse 更新API金鑰回應
// @Description 更新API金鑰後的回應
type UpdateAPIKeyResponse struct {
	BaseResponse
	Data *ExtendedAPIKeyVO `json:"data,omitempty"`
}

// ChangePasswordResponse 修改密碼回應
// @Description 修改密碼後的回應
type ChangePasswordResponse struct {
	BaseResponse
}

// ResetPasswordResponse 重設密碼回應
// @Description 重設密碼請求後的回應
type ResetPasswordResponse struct {
	BaseResponse
	Data *struct {
		Message string `json:"message" example:"Password reset email sent"`
	} `json:"data,omitempty"`
}

// ConfirmResetPasswordResponse 確認重設密碼回應
// @Description 確認重設密碼後的回應
type ConfirmResetPasswordResponse struct {
	BaseResponse
}

// LogoutResponse 登出回應
// @Description 登出後的回應
type LogoutResponse struct {
	BaseResponse
}

// DeleteUserResponse 刪除使用者回應
// @Description 刪除使用者後的回應
type DeleteUserResponse struct {
	BaseResponse
}

// DeleteAPIKeyResponse 刪除API金鑰回應
// @Description 刪除API金鑰後的回應
type DeleteAPIKeyResponse struct {
	BaseResponse
}

// ProfileStatsVO 個人檔案統計
type ProfileStatsVO struct {
	TotalAPIRequests    int64 `json:"total_api_requests" example:"15000"`
	MonthlyAPIRequests  int64 `json:"monthly_api_requests" example:"2500"`
	RemainingAPIQuota   int   `json:"remaining_api_quota" example:"7500"`
	TotalAPIKeys        int   `json:"total_api_keys" example:"3"`
	ActiveAPIKeys       int   `json:"active_api_keys" example:"2"`
	SubscriptionStatus  string `json:"subscription_status" example:"active"`
	DaysUntilExpiration *int   `json:"days_until_expiration,omitempty" example:"45"`
}

// GetProfileStatsResponse 取得個人檔案統計回應
// @Description 取得個人檔案統計資訊的回應
type GetProfileStatsResponse struct {
	BaseResponse
	Data *ProfileStatsVO `json:"data,omitempty"`
}

// AdminStatsVO 管理員統計
type AdminStatsVO struct {
	TotalUsers         int64 `json:"total_users" example:"1500"`
	ActiveUsers        int64 `json:"active_users" example:"1200"`
	NewUsersThisMonth  int64 `json:"new_users_this_month" example:"150"`
	TotalAPIRequests   int64 `json:"total_api_requests" example:"500000"`
	MonthlyAPIRequests int64 `json:"monthly_api_requests" example:"75000"`
	TotalAPIKeys       int64 `json:"total_api_keys" example:"4500"`
	ActiveAPIKeys      int64 `json:"active_api_keys" example:"3200"`
}

// GetAdminStatsResponse 取得管理員統計回應
// @Description 取得管理員統計資訊的回應
type GetAdminStatsResponse struct {
	BaseResponse
	Data *AdminStatsVO `json:"data,omitempty"`
}

// MaskAPIKey 遮蔽API金鑰敏感資訊
func (k *ExtendedAPIKeyVO) MaskAPIKey() {
	if k.Key != "" {
		keyLen := len(k.Key)
		if keyLen <= 12 {
			k.KeyPreview = k.Key[:4] + strings.Repeat("*", keyLen-4)
		} else {
			k.KeyPreview = k.Key[:8] + strings.Repeat("*", keyLen-12) + k.Key[keyLen-4:]
		}
		k.Key = "" // 清除完整金鑰
	}
}

// CalculateRemainingQuota 計算剩餘配額
func (u *UserVO) CalculateRemainingQuota() int {
	remaining := u.APIQuota - u.APIUsage
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsSubscriptionActive 檢查訂閱是否活躍
func (u *UserVO) IsSubscriptionActive() bool {
	if u.SubscriptionExpiresAt == nil {
		return true // 永久訂閱
	}
	return u.SubscriptionExpiresAt.After(time.Now())
}

// GetSubscriptionDaysRemaining 取得訂閱剩餘天數
func (u *UserVO) GetSubscriptionDaysRemaining() *int {
	if u.SubscriptionExpiresAt == nil {
		return nil // 永久訂閱
	}
	
	days := int(u.SubscriptionExpiresAt.Sub(time.Now()).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return &days
} 