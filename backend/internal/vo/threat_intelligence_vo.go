package vo

import (
	"time"

	"github.com/google/uuid"
)

// ThreatIntelligenceVO 威脅情報回應
type ThreatIntelligenceVO struct {
	ID              uuid.UUID              `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IPAddress       string                 `json:"ip_address" example:"192.168.1.100"`
	Domain          *string                `json:"domain" example:"malicious.example.com"`
	ThreatType      string                 `json:"threat_type" example:"malware" enums:"malware,phishing,spam,botnet,scanner,ddos,bruteforce,other"`
	Severity        string                 `json:"severity" example:"high" enums:"low,medium,high,critical"`
	ConfidenceScore int                    `json:"confidence_score" example:"85" minimum:"0" maximum:"100"`
	Description     *string                `json:"description" example:"Malicious IP detected by multiple sources"`
	Source          string                 `json:"source" example:"AbuseIPDB"`
	ExternalID      *string                `json:"external_id" example:"ext-12345"`
	CountryCode     *string                `json:"country_code" example:"TW"`
	ASN             *int                   `json:"asn" example:"8075"`
	ISP             *string                `json:"isp" example:"Microsoft Corporation"`
	FirstSeen       time.Time              `json:"first_seen" example:"2024-01-01T00:00:00Z"`
	LastSeen        time.Time              `json:"last_seen" example:"2024-01-02T00:00:00Z"`
	Tags            []string               `json:"tags" example:"botnet,malware"`
	Metadata        map[string]interface{} `json:"metadata" example:"{}"`
	RiskScore       int                    `json:"risk_score" example:"92" minimum:"0" maximum:"100"`
	IsHighRisk      bool                   `json:"is_high_risk" example:"true"`
	IsRecent        bool                   `json:"is_recent" example:"true"`
	CreatedAt       time.Time              `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       time.Time              `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ThreatIntelligenceListVO 威脅情報列表回應
type ThreatIntelligenceListVO struct {
	Data       []ThreatIntelligenceVO `json:"data"`
	Pagination PaginationVO           `json:"pagination"`
}

// ThreatIntelligenceStatsVO 威脅情報統計回應
type ThreatIntelligenceStatsVO struct {
	TotalThreats   int                    `json:"total_threats" example:"1000"`
	HighRiskCount  int                    `json:"high_risk_count" example:"250"`
	RecentCount    int                    `json:"recent_count" example:"50"`
	CountByType    map[string]int         `json:"count_by_type" example:"{\"malware\":300,\"phishing\":200}"`
	CountBySeverity map[string]int        `json:"count_by_severity" example:"{\"high\":250,\"medium\":500}"`
	CountBySource  map[string]int         `json:"count_by_source" example:"{\"AbuseIPDB\":600,\"Manual\":400}"`
	CountByCountry map[string]int         `json:"count_by_country" example:"{\"TW\":100,\"US\":200}"`
	Timeline       []ThreatTimelineVO     `json:"timeline"`
}

// ThreatTimelineVO 威脅時間線回應
type ThreatTimelineVO struct {
	Date  string `json:"date" example:"2024-01-01"`
	Count int    `json:"count" example:"25"`
}

// ThreatIntelligenceIPLookupVO IP 查詢回應
type ThreatIntelligenceIPLookupVO struct {
	IPAddress       string                 `json:"ip_address" example:"192.168.1.100"`
	IsKnownThreat   bool                   `json:"is_known_threat" example:"true"`
	ThreatCount     int                    `json:"threat_count" example:"3"`
	HighestSeverity string                 `json:"highest_severity" example:"high"`
	Sources         []string               `json:"sources" example:"AbuseIPDB,Manual"`
	FirstSeen       *time.Time             `json:"first_seen" example:"2024-01-01T00:00:00Z"`
	LastSeen        *time.Time             `json:"last_seen" example:"2024-01-02T00:00:00Z"`
	Details         []ThreatIntelligenceVO `json:"details"`
}

// ThreatIntelligenceDomainLookupVO 域名查詢回應
type ThreatIntelligenceDomainLookupVO struct {
	Domain          string                 `json:"domain" example:"malicious.example.com"`
	IsKnownThreat   bool                   `json:"is_known_threat" example:"true"`
	ThreatCount     int                    `json:"threat_count" example:"2"`
	HighestSeverity string                 `json:"highest_severity" example:"critical"`
	Sources         []string               `json:"sources" example:"AbuseIPDB,Manual"`
	FirstSeen       *time.Time             `json:"first_seen" example:"2024-01-01T00:00:00Z"`
	LastSeen        *time.Time             `json:"last_seen" example:"2024-01-02T00:00:00Z"`
	Details         []ThreatIntelligenceVO `json:"details"`
}

// ThreatIntelligenceBulkCreateVO 批量建立回應
type ThreatIntelligenceBulkCreateVO struct {
	Success      []ThreatIntelligenceVO `json:"success"`
	Failed       []BulkOperationError   `json:"failed"`
	TotalCount   int                    `json:"total_count" example:"100"`
	SuccessCount int                    `json:"success_count" example:"95"`
	FailedCount  int                    `json:"failed_count" example:"5"`
}

// ThreatIntelligenceBulkUpdateVO 批量更新回應
type ThreatIntelligenceBulkUpdateVO struct {
	Success      []ThreatIntelligenceVO `json:"success"`
	Failed       []BulkOperationError   `json:"failed"`
	TotalCount   int                    `json:"total_count" example:"50"`
	SuccessCount int                    `json:"success_count" example:"48"`
	FailedCount  int                    `json:"failed_count" example:"2"`
}

// ThreatIntelligenceBulkDeleteVO 批量刪除回應
type ThreatIntelligenceBulkDeleteVO struct {
	Success      []string             `json:"success"`
	Failed       []BulkOperationError `json:"failed"`
	TotalCount   int                  `json:"total_count" example:"50"`
	SuccessCount int                  `json:"success_count" example:"48"`
	FailedCount  int                  `json:"failed_count" example:"2"`
}

// BulkOperationError 批量操作錯誤
type BulkOperationError struct {
	Index   int    `json:"index" example:"5"`
	ID      string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Error   string `json:"error" example:"Invalid IP address"`
	Message string `json:"message" example:"The provided IP address is not valid"`
}

// ThreatIntelligenceExportVO 匯出回應
type ThreatIntelligenceExportVO struct {
	FileName  string `json:"file_name" example:"threats_2024-01-01.json"`
	FileSize  int64  `json:"file_size" example:"1024000"`
	Count     int    `json:"count" example:"1000"`
	Format    string `json:"format" example:"json" enums:"json,csv,xml"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// ThreatIntelligenceTagVO 標籤回應
type ThreatIntelligenceTagVO struct {
	Tags      []string `json:"tags" example:"botnet,malware"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ThreatIntelligenceMetricsVO 指標回應
type ThreatIntelligenceMetricsVO struct {
	TotalThreats        int                    `json:"total_threats" example:"10000"`
	ActiveThreats       int                    `json:"active_threats" example:"5000"`
	HighRiskThreats     int                    `json:"high_risk_threats" example:"1000"`
	RecentThreats       int                    `json:"recent_threats" example:"500"`
	ThreatsByType       map[string]int         `json:"threats_by_type"`
	ThreatsBySeverity   map[string]int         `json:"threats_by_severity"`
	ThreatsBySource     map[string]int         `json:"threats_by_source"`
	ThreatsByCountry    map[string]int         `json:"threats_by_country"`
	ThreatsByASN        map[string]int         `json:"threats_by_asn"`
	TopMaliciousIPs     []string               `json:"top_malicious_ips"`
	TopMaliciousDomains []string               `json:"top_malicious_domains"`
	RecentActivity      []ThreatTimelineVO     `json:"recent_activity"`
	CollectionStatus    map[string]interface{} `json:"collection_status"`
}

// ThreatIntelligenceHealthVO 健康狀態回應
type ThreatIntelligenceHealthVO struct {
	Status             string                 `json:"status" example:"healthy" enums:"healthy,degraded,unhealthy"`
	DatabaseStatus     string                 `json:"database_status" example:"connected"`
	CollectionStatus   string                 `json:"collection_status" example:"active"`
	LastUpdate         time.Time              `json:"last_update" example:"2024-01-01T00:00:00Z"`
	ThreatCount        int                    `json:"threat_count" example:"10000"`
	SourceCount        int                    `json:"source_count" example:"5"`
	ActiveSourceCount  int                    `json:"active_source_count" example:"4"`
	CollectionMetrics  map[string]interface{} `json:"collection_metrics"`
	SystemMetrics      map[string]interface{} `json:"system_metrics"`
} 