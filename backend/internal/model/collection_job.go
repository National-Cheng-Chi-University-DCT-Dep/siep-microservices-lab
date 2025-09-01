package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CollectionStatus 收集狀態
type CollectionStatus string

const (
	StatusPending    CollectionStatus = "pending"
	StatusInProgress CollectionStatus = "in_progress"
	StatusCompleted  CollectionStatus = "completed"
	StatusFailed     CollectionStatus = "failed"
)

// CollectionJob 收集任務模型
type CollectionJob struct {
	ID               uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SourceID         uuid.UUID        `gorm:"type:uuid;not null" json:"source_id"`
	Status           CollectionStatus `gorm:"type:collection_status;default:'pending'" json:"status"`
	StartedAt        *time.Time       `gorm:"column:started_at" json:"started_at"`
	CompletedAt      *time.Time       `gorm:"column:completed_at" json:"completed_at"`
	RecordsCollected int              `gorm:"default:0" json:"records_collected"`
	ErrorMessage     *string          `gorm:"type:text" json:"error_message"`
	CreatedAt        time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// 關聯
	Source IntelligenceSource `gorm:"foreignKey:SourceID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName 指定資料表名稱
func (CollectionJob) TableName() string {
	return "collection_jobs"
}

// BeforeCreate 在建立前執行
func (c *CollectionJob) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Start 開始收集任務
func (c *CollectionJob) Start() {
	c.Status = StatusInProgress
	now := time.Now()
	c.StartedAt = &now
}

// Complete 完成收集任務
func (c *CollectionJob) Complete(recordsCollected int) {
	c.Status = StatusCompleted
	c.RecordsCollected = recordsCollected
	now := time.Now()
	c.CompletedAt = &now
}

// Fail 任務失敗
func (c *CollectionJob) Fail(errorMessage string) {
	c.Status = StatusFailed
	c.ErrorMessage = &errorMessage
	now := time.Now()
	c.CompletedAt = &now
}

// IsRunning 檢查任務是否正在運行
func (c *CollectionJob) IsRunning() bool {
	return c.Status == StatusInProgress
}

// IsCompleted 檢查任務是否已完成
func (c *CollectionJob) IsCompleted() bool {
	return c.Status == StatusCompleted
}

// IsFailed 檢查任務是否失敗
func (c *CollectionJob) IsFailed() bool {
	return c.Status == StatusFailed
}

// GetDuration 取得任務執行時間
func (c *CollectionJob) GetDuration() *time.Duration {
	if c.StartedAt == nil {
		return nil
	}
	
	endTime := time.Now()
	if c.CompletedAt != nil {
		endTime = *c.CompletedAt
	}
	
	duration := endTime.Sub(*c.StartedAt)
	return &duration
} 