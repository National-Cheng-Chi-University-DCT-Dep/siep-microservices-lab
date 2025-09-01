package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey API 金鑰模型
type APIKey struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	KeyHash   string     `gorm:"type:varchar(255);not null" json:"-"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt *time.Time `gorm:"column:expires_at" json:"expires_at"`
	Quota     int        `gorm:"default:1000" json:"quota"`
	Usage     int        `gorm:"default:0" json:"usage"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	LastUsed  *time.Time `gorm:"column:last_used" json:"last_used"`

	// 關聯
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定資料表名稱
func (APIKey) TableName() string {
	return "api_keys"
}

// BeforeCreate 在建立前執行
func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// IsExpired 檢查金鑰是否過期
func (a *APIKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return a.ExpiresAt.Before(time.Now())
}

// IsValid 檢查金鑰是否有效
func (a *APIKey) IsValid() bool {
	return a.IsActive && !a.IsExpired()
}

// CanUse 檢查是否可以使用
func (a *APIKey) CanUse() bool {
	if !a.IsValid() {
		return false
	}
	return a.Usage < a.Quota
}

// IncrementUsage 增加使用次數
func (a *APIKey) IncrementUsage() {
	a.Usage++
	now := time.Now()
	a.LastUsed = &now
} 