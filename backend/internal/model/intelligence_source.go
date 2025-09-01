package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IntelligenceSource 情報來源模型
type IntelligenceSource struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name               string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	URL                *string   `gorm:"type:varchar(500)" json:"url"`
	APIKeyRequired     bool      `gorm:"default:false" json:"api_key_required"`
	IsActive           bool      `gorm:"default:true" json:"is_active"`
	CollectionInterval int       `gorm:"default:3600" json:"collection_interval"` // 秒
	LastCollection     *time.Time `gorm:"column:last_collection" json:"last_collection"`
	TotalCollected     int       `gorm:"default:0" json:"total_collected"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// 關聯
	CollectionJobs []CollectionJob `gorm:"foreignKey:SourceID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定資料表名稱
func (IntelligenceSource) TableName() string {
	return "intelligence_sources"
}

// BeforeCreate 在建立前執行
func (i *IntelligenceSource) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// ShouldCollect 檢查是否應該收集
func (i *IntelligenceSource) ShouldCollect() bool {
	if !i.IsActive {
		return false
	}
	
	if i.LastCollection == nil {
		return true
	}
	
	// 檢查是否已經過了收集間隔
	nextCollection := i.LastCollection.Add(time.Duration(i.CollectionInterval) * time.Second)
	return time.Now().After(nextCollection)
}

// UpdateLastCollection 更新最後收集時間
func (i *IntelligenceSource) UpdateLastCollection() {
	now := time.Now()
	i.LastCollection = &now
}

// IncrementCollected 增加收集數量
func (i *IntelligenceSource) IncrementCollected(count int) {
	i.TotalCollected += count
} 