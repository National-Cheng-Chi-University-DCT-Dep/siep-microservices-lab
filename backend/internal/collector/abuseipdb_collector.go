package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// AbuseIPDBCollector AbuseIPDB 威脅情報收集器
type AbuseIPDBCollector struct {
	apiKey  string
	baseURL string
	client  *http.Client
	service service.ThreatIntelligenceService
}

// AbuseIPDBData AbuseIPDB API 資料結構
type AbuseIPDBData struct {
	IPAddress            string `json:"ipAddress"`
	IsPublic             bool   `json:"isPublic"`
	IPVersion            int    `json:"ipVersion"`
	IsWhitelisted        bool   `json:"isWhitelisted"`
	AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
	CountryCode          string `json:"countryCode"`
	CountryName          string `json:"countryName"`
	UsageType            string `json:"usageType"`
	ISP                  string `json:"isp"`
	Domain               string `json:"domain"`
	TotalReports         int    `json:"totalReports"`
	NumDistinctUsers     int    `json:"numDistinctUsers"`
	LastReportedAt       string `json:"lastReportedAt"`
}

// AbuseIPDBResponse AbuseIPDB API 回應結構
type AbuseIPDBResponse struct {
	Data AbuseIPDBData `json:"data"`
}

// NewAbuseIPDBCollector 建立新的 AbuseIPDB 收集器
func NewAbuseIPDBCollector(apiKey string, service service.ThreatIntelligenceService) *AbuseIPDBCollector {
	return &AbuseIPDBCollector{
		apiKey:  apiKey,
		baseURL: "https://api.abuseipdb.com/api/v2",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		service: service,
	}
}

// CollectIPThreatIntel 收集指定 IP 的威脅情報
func (c *AbuseIPDBCollector) CollectIPThreatIntel(ctx context.Context, ipAddress string) error {
	// 檢查 IP 格式
	if !isValidIP(ipAddress) {
		return fmt.Errorf("無效的 IP 地址格式: %s", ipAddress)
	}

	// 呼叫 AbuseIPDB API
	response, err := c.queryAbuseIPDB(ipAddress)
	if err != nil {
		pkglogger.Error("AbuseIPDB API 查詢失敗", pkglogger.Fields{
			"ip_address": ipAddress,
			"error":      err.Error(),
		})
		return fmt.Errorf("AbuseIPDB API 查詢失敗: %w", err)
	}

	// 轉換為威脅情報記錄
	threatIntel := c.convertToThreatIntel(response)
	if threatIntel == nil {
		pkglogger.Info("IP 無威脅情報", pkglogger.Fields{
			"ip_address": ipAddress,
		})
		return nil
	}

	// 儲存到資料庫
	_, err = c.service.CreateThreat(ctx, threatIntel)
	if err != nil {
		pkglogger.Error("威脅情報儲存失敗", pkglogger.Fields{
			"ip_address": ipAddress,
			"error":      err.Error(),
		})
		return fmt.Errorf("威脅情報儲存失敗: %w", err)
	}

	pkglogger.Info("威脅情報收集成功", pkglogger.Fields{
		"ip_address":       ipAddress,
		"confidence_score": threatIntel.ConfidenceScore,
		"threat_type":      threatIntel.ThreatType,
	})

	return nil
}

// CollectBulkIPThreatIntel 批量收集 IP 威脅情報
func (c *AbuseIPDBCollector) CollectBulkIPThreatIntel(ctx context.Context, ipAddresses []string) ([]string, []string) {
	var successful []string
	var failed []string

	for _, ipAddress := range ipAddresses {
		select {
		case <-ctx.Done():
			// 檢查是否被取消
			pkglogger.Warn("批量收集被取消", pkglogger.Fields{
				"processed": len(successful) + len(failed),
				"total":     len(ipAddresses),
			})
			return successful, failed
		default:
			err := c.CollectIPThreatIntel(ctx, ipAddress)
			if err != nil {
				failed = append(failed, ipAddress)
				pkglogger.Warn("單一 IP 收集失敗", pkglogger.Fields{
					"ip_address": ipAddress,
					"error":      err.Error(),
				})
			} else {
				successful = append(successful, ipAddress)
			}

			// 避免 API 限流，每次查詢間隔 1 秒
			time.Sleep(1 * time.Second)
		}
	}

	pkglogger.Info("批量收集完成", pkglogger.Fields{
		"total":      len(ipAddresses),
		"successful": len(successful),
		"failed":     len(failed),
	})

	return successful, failed
}

// queryAbuseIPDB 查詢 AbuseIPDB API
func (c *AbuseIPDBCollector) queryAbuseIPDB(ipAddress string) (*AbuseIPDBResponse, error) {
	// 建立請求
	url := fmt.Sprintf("%s/check", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("建立 HTTP 請求失敗: %w", err)
	}

	// 設定查詢參數
	q := req.URL.Query()
	q.Add("ipAddress", ipAddress)
	q.Add("maxAgeInDays", "90") // 查詢 90 天內的記錄
	q.Add("verbose", "")        // 取得詳細資訊
	req.URL.RawQuery = q.Encode()

	// 設定請求標頭
	req.Header.Set("Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	// 發送請求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 請求失敗: %w", err)
	}
	defer resp.Body.Close()

	// 檢查狀態碼
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回錯誤狀態碼 %d: %s", resp.StatusCode, string(body))
	}

	// 解析回應
	var response AbuseIPDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析 JSON 回應失敗: %w", err)
	}

	return &response, nil
}

// convertToThreatIntel 轉換 AbuseIPDB 回應為威脅情報記錄
func (c *AbuseIPDBCollector) convertToThreatIntel(response *AbuseIPDBResponse) *dto.ThreatIntelligenceCreateRequest {
	data := response.Data

	// 如果信心分數太低，不建立記錄
	if data.AbuseConfidenceScore < 25 {
		return nil
	}

	// 判斷威脅類型
	threatType := c.determineThreatType(data.AbuseConfidenceScore, data.UsageType)
	
	// 判斷嚴重程度
	severity := c.determineSeverity(data.AbuseConfidenceScore)

	// 解析最後回報時間（暫時不使用，因為 DTO 沒有此欄位）
	// var lastSeen *time.Time
	// if data.LastReportedAt != "" {
	// 	if t, err := time.Parse(time.RFC3339, data.LastReportedAt); err == nil {
	// 		lastSeen = &t
	// 	}
	// }

	description := c.generateDescription(data)
	return &dto.ThreatIntelligenceCreateRequest{
		IPAddress:       data.IPAddress,
		ThreatType:      threatType,
		Severity:        severity,
		ConfidenceScore: data.AbuseConfidenceScore,
		Source:          "AbuseIPDB",
		Description:     &description,
		CountryCode:     &data.CountryCode,
		Tags:            c.generateTags(data),
		Metadata:        c.generateMetadata(data),
	}
}

// determineThreatType 根據信心分數和使用類型判斷威脅類型
func (c *AbuseIPDBCollector) determineThreatType(confidenceScore int, usageType string) string {
	switch {
	case confidenceScore >= 75:
		if strings.Contains(strings.ToLower(usageType), "hosting") {
			return "botnet"
		}
		return "scanner"
	case confidenceScore >= 50:
		return "spam"
	case confidenceScore >= 25:
		return "scanner"
	default:
		return "other"
	}
}

// determineSeverity 根據信心分數判斷嚴重程度
func (c *AbuseIPDBCollector) determineSeverity(confidenceScore int) string {
	switch {
	case confidenceScore >= 85:
		return "critical"
	case confidenceScore >= 70:
		return "high"
	case confidenceScore >= 50:
		return "medium"
	default:
		return "low"
	}
}

// generateDescription 生成威脅描述
func (c *AbuseIPDBCollector) generateDescription(data AbuseIPDBData) string {
	return fmt.Sprintf(
		"AbuseIPDB 回報的惡意 IP，信心分數: %d%%，總回報次數: %d，來源: %s (%s)",
		data.AbuseConfidenceScore,
		data.TotalReports,
		data.ISP,
		data.CountryName,
	)
}

// generateTags 生成標籤
func (c *AbuseIPDBCollector) generateTags(data AbuseIPDBData) []string {
	var tags []string
	
	tags = append(tags, "abuseipdb")
	
	if data.IsWhitelisted {
		tags = append(tags, "whitelisted")
	}
	
	if data.UsageType != "" {
		tags = append(tags, strings.ToLower(data.UsageType))
	}
	
	if data.CountryCode != "" {
		tags = append(tags, strings.ToLower(data.CountryCode))
	}

	return tags
}

// generateMetadata 生成元資料
func (c *AbuseIPDBCollector) generateMetadata(data AbuseIPDBData) map[string]interface{} {
	metadata := make(map[string]interface{})
	
	metadata["abuseipdb_confidence_score"] = data.AbuseConfidenceScore
	metadata["abuseipdb_total_reports"] = data.TotalReports
	metadata["abuseipdb_distinct_users"] = data.NumDistinctUsers
	metadata["ip_version"] = data.IPVersion
	metadata["is_public"] = data.IsPublic
	metadata["is_whitelisted"] = data.IsWhitelisted
	
	if data.ISP != "" {
		metadata["isp"] = data.ISP
	}
	
	if data.Domain != "" {
		metadata["domain"] = data.Domain
	}
	
	if data.UsageType != "" {
		metadata["usage_type"] = data.UsageType
	}

	return metadata
}

// isValidIP 簡單的 IP 地址格式檢查
func isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	
	for _, part := range parts {
		if num, err := strconv.Atoi(part); err != nil || num < 0 || num > 255 {
			return false
		}
	}
	
	return true
} 