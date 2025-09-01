package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ThreatType 威脅類型
type ThreatType string

const (
	ThreatMalware    ThreatType = "malware"
	ThreatPhishing   ThreatType = "phishing"
	ThreatSpam       ThreatType = "spam"
	ThreatBotnet     ThreatType = "botnet"
	ThreatScanner    ThreatType = "scanner"
	ThreatDDoS       ThreatType = "ddos"
	ThreatBruteforce ThreatType = "bruteforce"
	ThreatOther      ThreatType = "other"
)

// SeverityLevel 嚴重程度
type SeverityLevel string

const (
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)

// 使用共享的 JSONB 類型，定義在 quantum_job.go 中

// StringArray 自訂字串陣列類型
type StringArray []string

// Value 實作 driver.Valuer 介面
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	result := "{"
	for i, v := range s {
		if i > 0 {
			result += ","
		}
		result += `"` + v + `"`
	}
	result += "}"
	return result, nil
}

// Scan 實作 sql.Scanner 介面
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	// 簡單的 PostgreSQL 陣列解析
	str := string(bytes)
	if str == "{}" {
		*s = []string{}
		return nil
	}
	
	// 移除大括號並分割
	str = str[1 : len(str)-1]
	return json.Unmarshal([]byte("["+str+"]"), s)
}

// ThreatIntelligence 威脅情報模型
type ThreatIntelligence struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	IPAddress       net.IP        `gorm:"type:inet;not null" json:"ip_address"`
	Domain          *string       `gorm:"type:varchar(253)" json:"domain"`
	ThreatType      ThreatType    `gorm:"type:threat_type;not null" json:"threat_type"`
	Severity        SeverityLevel `gorm:"type:severity_level;not null" json:"severity"`
	ConfidenceScore int           `gorm:"check:confidence_score >= 0 AND confidence_score <= 100" json:"confidence_score"`
	Description     *string       `gorm:"type:text" json:"description"`
	Source          string        `gorm:"type:varchar(100);not null" json:"source"`
	ExternalID      *string       `gorm:"type:varchar(100)" json:"external_id"`
	CountryCode     *string       `gorm:"type:varchar(2)" json:"country_code"`
	ASN             *int          `gorm:"type:integer" json:"asn"`
	ISP             *string       `gorm:"type:varchar(200)" json:"isp"`
	FirstSeen       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"first_seen"`
	LastSeen        time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"last_seen"`
	Tags            StringArray   `gorm:"type:text[]" json:"tags"`
	Metadata        JSONB         `gorm:"type:jsonb" json:"metadata"`
	CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定資料表名稱
func (ThreatIntelligence) TableName() string {
	return "threat_intelligence"
}

// BeforeCreate 在建立前執行
func (t *ThreatIntelligence) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// IsHighRisk 檢查是否為高風險威脅
func (t *ThreatIntelligence) IsHighRisk() bool {
	return t.Severity == SeverityHigh || t.Severity == SeverityCritical
}

// IsRecent 檢查是否為最近的威脅（24小時內）
func (t *ThreatIntelligence) IsRecent() bool {
	return t.LastSeen.After(time.Now().Add(-24 * time.Hour))
}

// GetRiskScore 計算風險分數
func (t *ThreatIntelligence) GetRiskScore() int {
	baseScore := t.ConfidenceScore
	
	// 根據嚴重程度調整分數
	switch t.Severity {
	case SeverityCritical:
		baseScore = int(float64(baseScore) * 1.5)
	case SeverityHigh:
		baseScore = int(float64(baseScore) * 1.3)
	case SeverityMedium:
		baseScore = int(float64(baseScore) * 1.1)
	case SeverityLow:
		baseScore = int(float64(baseScore) * 0.9)
	}
	
	// 最近活動的威脅分數更高
	if t.IsRecent() {
		baseScore = int(float64(baseScore) * 1.2)
	}
	
	// 確保分數在 0-100 範圍內
	if baseScore > 100 {
		baseScore = 100
	}
	if baseScore < 0 {
		baseScore = 0
	}
	
	return baseScore
}

// AddTag 添加標籤
func (t *ThreatIntelligence) AddTag(tag string) {
	for _, existingTag := range t.Tags {
		if existingTag == tag {
			return // 標籤已存在
		}
	}
	t.Tags = append(t.Tags, tag)
}

// RemoveTag 移除標籤
func (t *ThreatIntelligence) RemoveTag(tag string) {
	for i, existingTag := range t.Tags {
		if existingTag == tag {
			t.Tags = append(t.Tags[:i], t.Tags[i+1:]...)
			return
		}
	}
}

// HasTag 檢查是否有特定標籤
func (t *ThreatIntelligence) HasTag(tag string) bool {
	for _, existingTag := range t.Tags {
		if existingTag == tag {
			return true
		}
	}
	return false
} 