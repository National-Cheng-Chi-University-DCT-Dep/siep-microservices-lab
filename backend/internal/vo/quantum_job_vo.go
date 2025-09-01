package vo

import (
	"time"

	"github.com/google/uuid"
)

// QuantumJobResponse 量子任務響應
type QuantumJobResponse struct {
	ID                  uuid.UUID              `json:"id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Title               string                 `json:"title" example:"DDoS攻擊模式分析"`
	Description         string                 `json:"description" example:"分析近期網路流量中的 DDoS 攻擊模式"`
	Status              string                 `json:"status" example:"completed"`
	Priority            int                    `json:"priority" example:"5"`
	CreatedAt           time.Time              `json:"created_at" example:"2025-11-15T12:30:45Z"`
	UpdatedAt           time.Time              `json:"updated_at" example:"2025-11-15T12:45:30Z"`
	StartedAt           *time.Time             `json:"started_at,omitempty" example:"2025-11-15T12:35:00Z"`
	CompletedAt         *time.Time             `json:"completed_at,omitempty" example:"2025-11-15T12:45:00Z"`
	ExecutionTimeSeconds int                   `json:"execution_time_seconds,omitempty" example:"600"`
	QuantumBackend      string                 `json:"quantum_backend,omitempty" example:"ibmq_qasm_simulator"`
	IsSimulation         bool                  `json:"is_simulation" example:"true"`
	Tags                []string               `json:"tags,omitempty" example:"['ddos', 'network', 'analysis']"`
	ConfidenceScore     float64                `json:"confidence_score,omitempty" example:"85.5"`
	IsMalicious         *bool                  `json:"is_malicious,omitempty" example:"true"`
	Source              string                 `json:"source,omitempty" example:"api"`
	InputParamsSummary  map[string]interface{} `json:"input_params_summary,omitempty"`
	ResultsSummary      map[string]interface{} `json:"results_summary,omitempty"`
}

// QuantumJobDetailResponse 量子任務詳細響應
type QuantumJobDetailResponse struct {
	QuantumJobResponse
	InputParams  map[string]interface{} `json:"input_params"`
	Results      map[string]interface{} `json:"results,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty" example:"量子設備連接超時"`
	Notes        string                 `json:"notes,omitempty" example:"這是一個分析特定 DDoS 攻擊模式的任務"`
	Logs         []QuantumJobLogResponse `json:"logs,omitempty"`
}

// QuantumJobLogResponse 量子任務日誌響應
type QuantumJobLogResponse struct {
	ID        uint      `json:"id" example:"1"`
	OldStatus string    `json:"old_status,omitempty" example:"pending"`
	NewStatus string    `json:"new_status" example:"running"`
	Message   string    `json:"message" example:"任務開始執行"`
	CreatedBy string    `json:"created_by" example:"system"`
	CreatedAt time.Time `json:"created_at" example:"2025-11-15T12:35:00Z"`
}

// SubmitQuantumJobResponse 提交量子任務響應
type SubmitQuantumJobResponse struct {
	JobID     uuid.UUID `json:"job_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Status    string    `json:"status" example:"pending"`
	Message   string    `json:"message" example:"任務已成功提交，等待執行"`
	CreatedAt time.Time `json:"created_at" example:"2025-11-15T12:30:45Z"`
}

// ListQuantumJobsResponse 列出量子任務響應
type ListQuantumJobsResponse struct {
	Jobs       []QuantumJobResponse `json:"jobs"`
	Pagination PaginationInfo      `json:"pagination"`
}

// ThreatAnalysisSubmitResponse 威脅分析提交響應
type ThreatAnalysisSubmitResponse struct {
	JobID       uuid.UUID `json:"job_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Status      string    `json:"status" example:"pending"`
	Title       string    `json:"title" example:"多源資料威脅分析"`
	Message     string    `json:"message" example:"威脅分析任務已提交"`
	EstimatedTime int     `json:"estimated_time_seconds" example:"300"`
}

// GetNextPendingJobResponse 獲取下一個待處理任務響應（內部使用）
type GetNextPendingJobResponse struct {
	JobID       uuid.UUID              `json:"job_id" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	InputParams map[string]interface{} `json:"input_params"`
	Title       string                 `json:"title" example:"DDoS攻擊模式分析"`
	Priority    int                    `json:"priority" example:"5"`
}

// BatchSubmitQuantumJobsResponse 批次提交量子任務響應
type BatchSubmitQuantumJobsResponse struct {
	SuccessCount int                       `json:"success_count" example:"8"`
	FailedCount  int                       `json:"failed_count" example:"2"`
	JobIDs       []uuid.UUID               `json:"job_ids"`
	Failures     []BatchSubmissionFailure  `json:"failures,omitempty"`
}

// BatchSubmissionFailure 批次提交失敗資訊
type BatchSubmissionFailure struct {
	Index   int    `json:"index" example:"3"`
	Message string `json:"message" example:"參數驗證失敗"`
}
