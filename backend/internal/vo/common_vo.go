package vo

import (
	"time"

	"github.com/google/uuid"
)

// BaseResponse 基礎回應結構
type BaseResponse struct {
	Success   bool        `json:"success" example:"true"`
	Message   string      `json:"message" example:"Operation successful"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorVO    `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	RequestID string      `json:"request_id,omitempty" example:"req-123456"`
}

// ErrorVO 錯誤回應結構
type ErrorVO struct {
	Code    string      `json:"code" example:"INVALID_INPUT"`
	Message string      `json:"message" example:"Invalid input parameters"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationVO 分頁回應結構
type PaginationVO struct {
	CurrentPage  int   `json:"current_page" example:"1"`
	PageSize     int   `json:"page_size" example:"20"`
	TotalPages   int   `json:"total_pages" example:"10"`
	TotalRecords int64 `json:"total_records" example:"200"`
	HasNext      bool  `json:"has_next" example:"true"`
	HasPrevious  bool  `json:"has_previous" example:"false"`
}

// HealthCheckVO 健康檢查回應
type HealthCheckVO struct {
	Status       string                 `json:"status" example:"healthy" enums:"healthy,degraded,unhealthy"`
	Version      string                 `json:"version" example:"1.0.0"`
	Timestamp    time.Time              `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Services     map[string]ServiceStatus `json:"services"`
	Database     DatabaseStatus         `json:"database"`
	SystemInfo   SystemInfo             `json:"system_info"`
}

// ServiceStatus 服務狀態
type ServiceStatus struct {
	Status      string            `json:"status" example:"healthy"`
	LastCheck   time.Time         `json:"last_check" example:"2024-01-01T00:00:00Z"`
	ResponseTime string           `json:"response_time" example:"10ms"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// DatabaseStatus 資料庫狀態
type DatabaseStatus struct {
	Status         string `json:"status" example:"connected"`
	ConnectionTime string `json:"connection_time" example:"5ms"`
	OpenConnections int   `json:"open_connections" example:"10"`
	IdleConnections int   `json:"idle_connections" example:"5"`
	MaxConnections  int   `json:"max_connections" example:"25"`
}

// SystemInfo 系統資訊
type SystemInfo struct {
	GoVersion  string `json:"go_version" example:"go1.23.0"`
	Platform   string `json:"platform" example:"darwin/arm64"`
	StartTime  time.Time `json:"start_time" example:"2024-01-01T00:00:00Z"`
	Uptime     string `json:"uptime" example:"2h30m"`
	MemoryUsage string `json:"memory_usage" example:"45MB"`
	CPUUsage   string `json:"cpu_usage" example:"12%"`
}

// UserVO 使用者回應
type UserVO struct {
	ID                    uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username              string     `json:"username" example:"admin"`
	Email                 string     `json:"email" example:"admin@example.com"`
	Role                  string     `json:"role" example:"admin" enums:"admin,premium,basic"`
	IsActive              bool       `json:"is_active" example:"true"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at" example:"2024-12-31T23:59:59Z"`
	APIQuota              int        `json:"api_quota" example:"10000"`
	APIUsage              int        `json:"api_usage" example:"1250"`
	CreatedAt             time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt             time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	LastLogin             *time.Time `json:"last_login" example:"2024-01-01T12:00:00Z"`
}

// APIKeyVO API 金鑰回應
type APIKeyVO struct {
	ID        uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string     `json:"name" example:"My API Key"`
	IsActive  bool       `json:"is_active" example:"true"`
	ExpiresAt *time.Time `json:"expires_at" example:"2024-12-31T23:59:59Z"`
	Quota     int        `json:"quota" example:"1000"`
	Usage     int        `json:"usage" example:"150"`
	CreatedAt time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	LastUsed  *time.Time `json:"last_used" example:"2024-01-01T10:00:00Z"`
}

// IntelligenceSourceVO 情報來源回應
type IntelligenceSourceVO struct {
	ID                 uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name               string     `json:"name" example:"AbuseIPDB"`
	URL                *string    `json:"url" example:"https://api.abuseipdb.com/api/v2"`
	APIKeyRequired     bool       `json:"api_key_required" example:"true"`
	IsActive           bool       `json:"is_active" example:"true"`
	CollectionInterval int        `json:"collection_interval" example:"3600"`
	LastCollection     *time.Time `json:"last_collection" example:"2024-01-01T10:00:00Z"`
	TotalCollected     int        `json:"total_collected" example:"50000"`
	CreatedAt          time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt          time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// CollectionJobVO 收集任務回應
type CollectionJobVO struct {
	ID               uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	SourceID         uuid.UUID  `json:"source_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	SourceName       string     `json:"source_name" example:"AbuseIPDB"`
	Status           string     `json:"status" example:"completed" enums:"pending,in_progress,completed,failed"`
	StartedAt        *time.Time `json:"started_at" example:"2024-01-01T10:00:00Z"`
	CompletedAt      *time.Time `json:"completed_at" example:"2024-01-01T10:05:00Z"`
	RecordsCollected int        `json:"records_collected" example:"250"`
	ErrorMessage     *string    `json:"error_message" example:"Connection timeout"`
	Duration         *string    `json:"duration" example:"5m30s"`
	CreatedAt        time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt        time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// StatsVO 統計回應
type StatsVO struct {
	Label string `json:"label" example:"Malware"`
	Value int    `json:"value" example:"1500"`
	Percentage float64 `json:"percentage" example:"15.5"`
}

// DashboardVO 儀表板回應
type DashboardVO struct {
	Overview      OverviewVO             `json:"overview"`
	RecentThreats []ThreatIntelligenceVO `json:"recent_threats"`
	RecentJobs    []CollectionJobVO      `json:"recent_jobs"`
	Statistics    []StatsVO              `json:"statistics"`
	Timeline      []ThreatTimelineVO     `json:"timeline"`
	Alerts        []AlertVO              `json:"alerts"`
}

// OverviewVO 概覽回應
type OverviewVO struct {
	TotalThreats      int `json:"total_threats" example:"10000"`
	ActiveThreats     int `json:"active_threats" example:"5000"`
	HighRiskThreats   int `json:"high_risk_threats" example:"1000"`
	RecentThreats     int `json:"recent_threats" example:"500"`
	TotalSources      int `json:"total_sources" example:"5"`
	ActiveSources     int `json:"active_sources" example:"4"`
	PendingJobs       int `json:"pending_jobs" example:"2"`
	CompletedJobs     int `json:"completed_jobs" example:"150"`
}

// AlertVO 警報回應
type AlertVO struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Level       string    `json:"level" example:"warning" enums:"info,warning,error,critical"`
	Title       string    `json:"title" example:"High Risk Threat Detected"`
	Message     string    `json:"message" example:"Multiple high-risk threats detected from same IP"`
	Source      string    `json:"source" example:"threat_analyzer"`
	IsRead      bool      `json:"is_read" example:"false"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
} 