package dto

import (
	"net"
	"time"
)

// ThreatIntelligenceCreateRequest 建立威脅情報請求
type ThreatIntelligenceCreateRequest struct {
	IPAddress       string                 `json:"ip_address" binding:"required,ip" validate:"required,ip"`
	Domain          *string                `json:"domain" validate:"omitempty,fqdn"`
	ThreatType      string                 `json:"threat_type" binding:"required,oneof=malware phishing spam botnet scanner ddos bruteforce other"`
	Severity        string                 `json:"severity" binding:"required,oneof=low medium high critical"`
	ConfidenceScore int                    `json:"confidence_score" binding:"required,min=0,max=100"`
	Description     *string                `json:"description" validate:"omitempty,max=1000"`
	Source          string                 `json:"source" binding:"required,min=1,max=100"`
	ExternalID      *string                `json:"external_id" validate:"omitempty,max=100"`
	CountryCode     *string                `json:"country_code" validate:"omitempty,len=2"`
	ASN             *int                   `json:"asn" validate:"omitempty,min=1"`
	ISP             *string                `json:"isp" validate:"omitempty,max=200"`
	Tags            []string               `json:"tags" validate:"omitempty"`
	Metadata        map[string]interface{} `json:"metadata" validate:"omitempty"`
}

// ThreatIntelligenceUpdateRequest 更新威脅情報請求
type ThreatIntelligenceUpdateRequest struct {
	Domain          *string                `json:"domain" validate:"omitempty,fqdn"`
	ThreatType      *string                `json:"threat_type" validate:"omitempty,oneof=malware phishing spam botnet scanner ddos bruteforce other"`
	Severity        *string                `json:"severity" validate:"omitempty,oneof=low medium high critical"`
	ConfidenceScore *int                   `json:"confidence_score" validate:"omitempty,min=0,max=100"`
	Description     *string                `json:"description" validate:"omitempty,max=1000"`
	CountryCode     *string                `json:"country_code" validate:"omitempty,len=2"`
	ASN             *int                   `json:"asn" validate:"omitempty,min=1"`
	ISP             *string                `json:"isp" validate:"omitempty,max=200"`
	Tags            []string               `json:"tags" validate:"omitempty"`
	Metadata        map[string]interface{} `json:"metadata" validate:"omitempty"`
}

// ThreatIntelligenceQueryRequest 查詢威脅情報請求
type ThreatIntelligenceQueryRequest struct {
	IPAddress   *string  `json:"ip_address" form:"ip_address" validate:"omitempty,ip"`
	Domain      *string  `json:"domain" form:"domain" validate:"omitempty,fqdn"`
	ThreatType  *string  `json:"threat_type" form:"threat_type" validate:"omitempty,oneof=malware phishing spam botnet scanner ddos bruteforce other"`
	Severity    *string  `json:"severity" form:"severity" validate:"omitempty,oneof=low medium high critical"`
	Source      *string  `json:"source" form:"source" validate:"omitempty,max=100"`
	CountryCode *string  `json:"country_code" form:"country_code" validate:"omitempty,len=2"`
	Tags        []string `json:"tags" form:"tags" validate:"omitempty"`
	
	// 分頁參數
	Page     int `json:"page" form:"page" validate:"omitempty,min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"omitempty,min=1,max=100"`
	
	// 排序參數
	SortBy    string `json:"sort_by" form:"sort_by" validate:"omitempty,oneof=created_at updated_at last_seen confidence_score"`
	SortOrder string `json:"sort_order" form:"sort_order" validate:"omitempty,oneof=asc desc"`
	
	// 時間範圍
	StartTime *time.Time `json:"start_time" form:"start_time" validate:"omitempty"`
	EndTime   *time.Time `json:"end_time" form:"end_time" validate:"omitempty"`
}

// ThreatIntelligenceBulkCreateRequest 批量建立威脅情報請求
type ThreatIntelligenceBulkCreateRequest struct {
	Items []ThreatIntelligenceCreateRequest `json:"items" binding:"required,min=1,max=100"`
}

// ThreatIntelligenceBulkUpdateRequest 批量更新威脅情報請求
type ThreatIntelligenceBulkUpdateRequest struct {
	Items []struct {
		ID   string                          `json:"id" binding:"required,uuid"`
		Data ThreatIntelligenceUpdateRequest `json:"data" binding:"required"`
	} `json:"items" binding:"required,min=1,max=100"`
}

// ThreatIntelligenceBulkDeleteRequest 批量刪除威脅情報請求
type ThreatIntelligenceBulkDeleteRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,max=100"`
}

// ThreatIntelligenceTagRequest 標籤操作請求
type ThreatIntelligenceTagRequest struct {
	Tags []string `json:"tags" binding:"required,min=1"`
}

// ThreatIntelligenceIPLookupRequest IP 查詢請求
type ThreatIntelligenceIPLookupRequest struct {
	IPAddress string `json:"ip_address" form:"ip_address" binding:"required,ip"`
}

// ThreatIntelligenceDomainLookupRequest 域名查詢請求
type ThreatIntelligenceDomainLookupRequest struct {
	Domain string `json:"domain" form:"domain" binding:"required,fqdn"`
}

// ThreatIntelligenceStatsRequest 統計請求
type ThreatIntelligenceStatsRequest struct {
	GroupBy   string     `json:"group_by" form:"group_by" validate:"omitempty,oneof=threat_type severity source country_code"`
	StartTime *time.Time `json:"start_time" form:"start_time" validate:"omitempty"`
	EndTime   *time.Time `json:"end_time" form:"end_time" validate:"omitempty"`
}

// ValidateIPAddress 驗證 IP 地址
func (r *ThreatIntelligenceCreateRequest) ValidateIPAddress() error {
	if net.ParseIP(r.IPAddress) == nil {
		return ErrInvalidIPAddress
	}
	return nil
}

// SetDefaults 設定預設值
func (r *ThreatIntelligenceQueryRequest) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
	if r.SortBy == "" {
		r.SortBy = "created_at"
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}

// GetOffset 取得分頁偏移量
func (r *ThreatIntelligenceQueryRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// GetLimit 取得限制數量
func (r *ThreatIntelligenceQueryRequest) GetLimit() int {
	return r.PageSize
} 