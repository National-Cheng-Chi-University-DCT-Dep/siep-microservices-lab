package dto

import (
	"github.com/google/uuid"
)

// SubmitQuantumJobRequest 提交量子任務請求
type SubmitQuantumJobRequest struct {
	Title       string                 `json:"title" binding:"required,min=3,max=100" example:"DDoS攻擊模式分析"`
	Description string                 `json:"description" binding:"max=500" example:"分析近期網路流量中的 DDoS 攻擊模式"`
	Priority    int                    `json:"priority" binding:"min=1,max=10" example:"5"`
	InputParams map[string]interface{} `json:"input_params" binding:"required"`
	Tags        []string               `json:"tags" binding:"omitempty,dive,max=30"`
	Notes       string                 `json:"notes" binding:"max=1000" example:"這是一個分析特定 DDoS 攻擊模式的任務"`
	Source      string                 `json:"source" binding:"omitempty,max=50" example:"api"`
}

// GetQuantumJobRequest 獲取量子任務請求
type GetQuantumJobRequest struct {
	JobID uuid.UUID `uri:"job_id" binding:"required,uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
}

// ListQuantumJobsRequest 列出量子任務請求
type ListQuantumJobsRequest struct {
	Status    string `form:"status" binding:"omitempty,oneof=pending running completed failed" example:"completed"`
	Page      int    `form:"page" binding:"min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"min=1,max=100" example:"20"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at priority title status" example:"created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// UpdateQuantumJobRequest 更新量子任務請求
type UpdateQuantumJobRequest struct {
	JobID       uuid.UUID              `uri:"job_id" binding:"required,uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Title       string                 `json:"title" binding:"omitempty,min=3,max=100" example:"更新後的任務標題"`
	Description string                 `json:"description" binding:"omitempty,max=500" example:"更新後的任務描述"`
	Priority    int                    `json:"priority" binding:"omitempty,min=1,max=10" example:"8"`
	Notes       string                 `json:"notes" binding:"omitempty,max=1000" example:"更新後的筆記"`
	Tags        []string               `json:"tags" binding:"omitempty,dive,max=30"`
	InputParams map[string]interface{} `json:"input_params" binding:"omitempty"`
}

// CancelQuantumJobRequest 取消量子任務請求
type CancelQuantumJobRequest struct {
	JobID  uuid.UUID `uri:"job_id" binding:"required,uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Reason string    `json:"reason" binding:"omitempty,max=500" example:"任務不再需要"`
}

// CompleteQuantumJobRequest 完成量子任務請求（內部使用）
type CompleteQuantumJobRequest struct {
	JobID           uuid.UUID              `json:"job_id" binding:"required,uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Results         map[string]interface{} `json:"results" binding:"required"`
	ConfidenceScore float64                `json:"confidence_score" binding:"required,min=0,max=100" example:"85.5"`
	IsMalicious     bool                   `json:"is_malicious" example:"true"`
	ExecutionTime   int                    `json:"execution_time_seconds" binding:"min=0" example:"120"`
	QuantumBackend  string                 `json:"quantum_backend" binding:"required" example:"ibmq_qasm_simulator"`
}

// FailQuantumJobRequest 標記任務失敗請求（內部使用）
type FailQuantumJobRequest struct {
	JobID        uuid.UUID `json:"job_id" binding:"required,uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	ErrorMessage string    `json:"error_message" binding:"required" example:"量子設備連接超時"`
	ExecutionTime int      `json:"execution_time_seconds" binding:"min=0" example:"45"`
}

// GetNextPendingJobRequest 獲取下一個待處理任務請求（內部使用）
type GetNextPendingJobRequest struct {
	WorkerID string `json:"worker_id" binding:"required" example:"quantum-worker-1"`
	Secret   string `json:"secret" binding:"required" example:"your-secret-key"`
}

// BatchSubmitQuantumJobsRequest 批次提交量子任務請求
type BatchSubmitQuantumJobsRequest struct {
	Jobs []SubmitQuantumJobRequest `json:"jobs" binding:"required,min=1,max=10,dive"`
}

// ThreatAnalysisRequest 威脅分析請求（特定場景）
type ThreatAnalysisRequest struct {
	Title        string   `json:"title" binding:"required" example:"多源資料威脅分析"`
	Description  string   `json:"description" binding:"omitempty" example:"分析來自多個來源的威脅情報"`
	Priority     int      `json:"priority" binding:"omitempty,min=1,max=10" example:"7"`
	DataSources  []string `json:"data_sources" binding:"required,min=1" example:"['abuseipdb', 'hibp', 'internal']"`
	ThreatType   string   `json:"threat_type" binding:"required" example:"malware"`
	TimeWindow   string   `json:"time_window" binding:"required" example:"24h"`
	UseSimulator bool     `json:"use_simulator" example:"true"`
}
