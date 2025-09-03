package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole 使用者角色
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RolePremium UserRole = "premium"
	RoleBasic   UserRole = "basic"
)

// User 使用者模型
type User struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username              string     `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email                 string     `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash          string     `gorm:"type:varchar(255);not null" json:"-"`
	Role                  UserRole   `gorm:"type:user_role;default:'basic'" json:"role"`
	IsActive              bool       `gorm:"default:true" json:"is_active"`
	EmailVerified         bool       `gorm:"default:false" json:"email_verified"`
	SubscriptionType      string     `gorm:"type:varchar(20);default:'basic'" json:"subscription_type"`
	SubscriptionExpiresAt *time.Time `gorm:"column:subscription_expires_at" json:"subscription_expires_at"`
	APIQuota              int        `gorm:"default:1000" json:"api_quota"`
	APIUsage              int        `gorm:"default:0" json:"api_usage"`
	UsedAPIQuota          int        `gorm:"default:0" json:"used_api_quota"`
	CreatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	LastLogin             *time.Time `gorm:"column:last_login" json:"last_login"`

	// 關聯
	APIKeys []APIKey `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定資料表名稱
func (User) TableName() string {
	return "users"
}

// BeforeCreate 在建立前執行
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// IsSubscriptionActive 檢查訂閱是否有效
func (u *User) IsSubscriptionActive() bool {
	if u.SubscriptionExpiresAt == nil {
		return false
	}
	return u.SubscriptionExpiresAt.After(time.Now())
}

// CanAccessAPI 檢查是否可以存取 API
func (u *User) CanAccessAPI() bool {
	if !u.IsActive {
		return false
	}
	
	// 管理員無限制
	if u.Role == RoleAdmin {
		return true
	}
	
	// 檢查 API 配額
	return u.APIUsage < u.APIQuota
}

// IncrementAPIUsage 增加 API 使用次數
func (u *User) IncrementAPIUsage() {
	u.APIUsage++
}

// ResetAPIUsage 重置 API 使用次數
func (u *User) ResetAPIUsage() {
	u.APIUsage = 0
} 