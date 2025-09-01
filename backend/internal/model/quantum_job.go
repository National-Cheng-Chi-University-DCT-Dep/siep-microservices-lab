package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QuantumJobStatus 量子任務狀態類型
type QuantumJobStatus string

const (
	JobStatusPending   QuantumJobStatus = "pending"   // 待處理
	JobStatusRunning   QuantumJobStatus = "running"   // 執行中
	JobStatusCompleted QuantumJobStatus = "completed" // 已完成
	JobStatusFailed    QuantumJobStatus = "failed"    // 失敗
)

// 使用共享的 JSONB 類型，定義在 common.go 中

// QuantumJob 量子任務模型
type QuantumJob struct {
	// 基本資訊
	ID          uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	Title       string           `gorm:"type:varchar(100);not null" json:"title"`
	Description string           `gorm:"type:text" json:"description"`
	Status      QuantumJobStatus `gorm:"type:quantum_job_status;not null;default:'pending'" json:"status"`
	Priority    int              `gorm:"not null;default:1" json:"priority"`

	// 時間戳
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`

	// 參數與結果
	InputParams   JSONB  `gorm:"type:jsonb;not null" json:"input_params"`
	Results       JSONB  `gorm:"type:jsonb" json:"results"`
	ErrorMessage  string `gorm:"type:text" json:"error_message"`

	// 執行資訊
	ExecutionTimeSeconds int     `json:"execution_time_seconds"`
	QuantumBackend      string  `gorm:"type:varchar(100)" json:"quantum_backend"`
	IsSimulation         bool    `gorm:"default:true" json:"is_simulation"`
	Shots                int     `gorm:"default:1024" json:"shots"`

	// 統計與分析資訊
	ConfidenceScore float64 `gorm:"type:decimal(5,2)" json:"confidence_score"`
	IsMalicious     *bool   `json:"is_malicious"`

	// 額外資訊
	Tags   []string `gorm:"type:text[]" json:"tags"`
	Notes  string   `gorm:"type:text" json:"notes"`
	Source string   `gorm:"type:varchar(50)" json:"source"`

	// 關聯
	User User `gorm:"foreignkey:UserID" json:"-"`
	Logs []QuantumJobLog `gorm:"foreignKey:JobID" json:"logs,omitempty"`
}

// TableName 指定資料表名稱
func (QuantumJob) TableName() string {
	return "quantum_jobs"
}

// BeforeCreate GORM 鉤子，在創建前執行
func (q *QuantumJob) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

// QuantumJobLog 量子任務日誌模型
type QuantumJobLog struct {
	ID        uint             `gorm:"primary_key;autoIncrement" json:"id"`
	JobID     uuid.UUID        `gorm:"type:uuid;not null" json:"job_id"`
	OldStatus QuantumJobStatus `gorm:"type:quantum_job_status" json:"old_status"`
	NewStatus QuantumJobStatus `gorm:"type:quantum_job_status;not null" json:"new_status"`
	Message   string           `gorm:"type:text" json:"message"`
	CreatedBy string           `gorm:"type:varchar(50);not null" json:"created_by"`
	CreatedAt time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 關聯
	Job QuantumJob `gorm:"foreignkey:JobID" json:"-"`
}

// TableName 指定資料表名稱
func (QuantumJobLog) TableName() string {
	return "quantum_job_logs"
}

// NewQuantumJob 創建新的量子任務
func NewQuantumJob(userID uuid.UUID, title string, inputParams map[string]interface{}) *QuantumJob {
	return &QuantumJob{
		UserID:      userID,
		Title:       title,
		Status:      JobStatusPending,
		Priority:    1,
		InputParams: inputParams,
		IsSimulation: true,
		Shots:       1024,
		Source:      "api",
	}
}

// MarkAsRunning 將任務標記為執行中
func (q *QuantumJob) MarkAsRunning() {
	q.Status = JobStatusRunning
	now := time.Now()
	q.StartedAt = &now
}

// MarkAsCompleted 將任務標記為已完成
func (q *QuantumJob) MarkAsCompleted(results map[string]interface{}, confidenceScore float64, isMalicious bool) {
	q.Status = JobStatusCompleted
	now := time.Now()
	q.CompletedAt = &now
	q.Results = results
	q.ConfidenceScore = confidenceScore
	q.IsMalicious = &isMalicious

	if q.StartedAt != nil {
		q.ExecutionTimeSeconds = int(now.Sub(*q.StartedAt).Seconds())
	}
}

// MarkAsFailed 將任務標記為失敗
func (q *QuantumJob) MarkAsFailed(errorMessage string) {
	q.Status = JobStatusFailed
	now := time.Now()
	q.CompletedAt = &now
	q.ErrorMessage = errorMessage

	if q.StartedAt != nil {
		q.ExecutionTimeSeconds = int(now.Sub(*q.StartedAt).Seconds())
	}
}

// IsFinished 檢查任務是否已完成（成功或失敗）
func (q *QuantumJob) IsFinished() bool {
	return q.Status == JobStatusCompleted || q.Status == JobStatusFailed
}
