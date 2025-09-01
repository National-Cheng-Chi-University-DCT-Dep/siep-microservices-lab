package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"crypto/sha1"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
)

// HIBPCollector 整合 Have I Been Pwned API v3
type HIBPCollector struct {
	apiKey     string
	baseURL    string
	client     *http.Client
	service    service.ThreatIntelligenceService
	userAgent  string
}

// HIBPBreach 數據泄露事件模型
type HIBPBreach struct {
	Name              string    `json:"Name"`
	Title             string    `json:"Title"`
	Domain            string    `json:"Domain"`
	BreachDate        string    `json:"BreachDate"`
	AddedDate         time.Time `json:"AddedDate"`
	ModifiedDate      time.Time `json:"ModifiedDate"`
	PwnCount          int       `json:"PwnCount"`
	Description       string    `json:"Description"`
	LogoPath          string    `json:"LogoPath"`
	DataClasses       []string  `json:"DataClasses"`
	IsVerified        bool      `json:"IsVerified"`
	IsFabricated      bool      `json:"IsFabricated"`
	IsSensitive       bool      `json:"IsSensitive"`
	IsRetired         bool      `json:"IsRetired"`
	IsSpamList        bool      `json:"IsSpamList"`
	IsMalware         bool      `json:"IsMalware"`
	IsStealerLog      bool      `json:"IsStealerLog"`
	IsSubscriptionFree bool     `json:"IsSubscriptionFree"`
}

// HIBPPaste Paste 模型
type HIBPPaste struct {
	Source     string    `json:"Source"`
	Id         string    `json:"Id"`
	Title      string    `json:"Title"`
	Date       time.Time `json:"Date"`
	EmailCount int       `json:"EmailCount"`
}

// HIBPSubscriptionStatus 訂閱狀態
type HIBPSubscriptionStatus struct {
	Name                    string `json:"Name"`
	Description             string `json:"Description"`
	Price                   int    `json:"Price"`
	PwnCount                int    `json:"PwnCount"`
	BreachCount             int    `json:"BreachCount"`
	PasteCount              int    `json:"PasteCount"`
	DomainSearchEnabled     bool   `json:"DomainSearchEnabled"`
	DomainSearchCapacity    int    `json:"DomainSearchCapacity"`
	DomainSearchUsed        int    `json:"DomainSearchUsed"`
	StealerLogEnabled       bool   `json:"StealerLogEnabled"`
	StealerLogCapacity      int    `json:"StealerLogCapacity"`
	StealerLogUsed          int    `json:"StealerLogUsed"`
	NextRenewalDate         string `json:"NextRenewalDate"`
	NextRenewalDateUtc      string `json:"NextRenewalDateUtc"`
	NextRenewalDateLocal    string `json:"NextRenewalDateLocal"`
	NextRenewalDateTimezone string `json:"NextRenewalDateTimezone"`
}

// HIBPDomainBreach 域名泄露查詢結果
type HIBPDomainBreach map[string][]string

// NewHIBPCollector 創建 HIBP 收集器
func NewHIBPCollector(apiKey string, service service.ThreatIntelligenceService) *HIBPCollector {
	return &HIBPCollector{
		apiKey:    apiKey,
		baseURL:   "https://haveibeenpwned.com/api/v3",
		userAgent: "Security-Intelligence-Platform/1.0",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		service: service,
	}
}

// CheckAccountBreaches 檢查帳戶是否在泄露事件中
func (c *HIBPCollector) CheckAccountBreaches(ctx context.Context, account string, includeUnverified bool) ([]HIBPBreach, error) {
	endpoint := fmt.Sprintf("/breachedaccount/%s", url.QueryEscape(account))
	
	// 添加查詢參數
	params := url.Values{}
	if !includeUnverified {
		params.Set("IncludeUnverified", "false")
	}
	
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	var breaches []HIBPBreach
	err := c.makeRequest(ctx, endpoint, &breaches)
	return breaches, err
}

// GetDomainBreaches 獲取域名下的所有泄露郵箱
func (c *HIBPCollector) GetDomainBreaches(ctx context.Context, domain string) (HIBPDomainBreach, error) {
	endpoint := fmt.Sprintf("/breacheddomain/%s", url.QueryEscape(domain))
	
	var domainBreaches HIBPDomainBreach
	err := c.makeRequest(ctx, endpoint, &domainBreaches)
	return domainBreaches, err
}

// GetSubscribedDomains 獲取已訂閱的域名列表
func (c *HIBPCollector) GetSubscribedDomains(ctx context.Context) ([]map[string]interface{}, error) {
	var domains []map[string]interface{}
	err := c.makeRequest(ctx, "/subscribeddomains", &domains)
	return domains, err
}

// GetAllBreaches 獲取所有泄露事件
func (c *HIBPCollector) GetAllBreaches(ctx context.Context, domain string) ([]HIBPBreach, error) {
	endpoint := "/breaches"
	
	if domain != "" {
		endpoint += "?domain=" + url.QueryEscape(domain)
	}
	
	var breaches []HIBPBreach
	err := c.makeRequest(ctx, endpoint, &breaches)
	return breaches, err
}

// GetBreachByName 根據名稱獲取特定泄露事件
func (c *HIBPCollector) GetBreachByName(ctx context.Context, name string) (*HIBPBreach, error) {
	endpoint := fmt.Sprintf("/breach/%s", url.QueryEscape(name))
	
	var breach HIBPBreach
	err := c.makeRequest(ctx, endpoint, &breach)
	return &breach, err
}

// GetLatestBreach 獲取最新添加的泄露事件
func (c *HIBPCollector) GetLatestBreach(ctx context.Context) (*HIBPBreach, error) {
	var breach HIBPBreach
	err := c.makeRequest(ctx, "/latestbreach", &breach)
	return &breach, err
}

// GetDataClasses 獲取所有數據類別
func (c *HIBPCollector) GetDataClasses(ctx context.Context) ([]string, error) {
	var dataClasses []string
	err := c.makeRequest(ctx, "/dataclasses", &dataClasses)
	return dataClasses, err
}

// GetAccountPastes 獲取帳戶的所有 Paste 記錄
func (c *HIBPCollector) GetAccountPastes(ctx context.Context, account string) ([]HIBPPaste, error) {
	endpoint := fmt.Sprintf("/pasteaccount/%s", url.QueryEscape(account))
	
	var pastes []HIBPPaste
	err := c.makeRequest(ctx, endpoint, &pastes)
	return pastes, err
}

// GetSubscriptionStatus 獲取訂閱狀態
func (c *HIBPCollector) GetSubscriptionStatus(ctx context.Context) (*HIBPSubscriptionStatus, error) {
	var status HIBPSubscriptionStatus
	err := c.makeRequest(ctx, "/subscription/status", &status)
	return &status, err
}

// Stealer Logs 相關方法 (需要 Pwned 5 或更高訂閱)
func (c *HIBPCollector) GetStealerLogsByEmail(ctx context.Context, email string) ([]string, error) {
	endpoint := fmt.Sprintf("/stealerlogsbyemail/%s", url.QueryEscape(email))
	
	var domains []string
	err := c.makeRequest(ctx, endpoint, &domains)
	return domains, err
}

func (c *HIBPCollector) GetStealerLogsByWebsiteDomain(ctx context.Context, domain string) ([]string, error) {
	endpoint := fmt.Sprintf("/stealerlogsbywebsitedomain/%s", url.QueryEscape(domain))
	
	var emails []string
	err := c.makeRequest(ctx, endpoint, &emails)
	return emails, err
}

func (c *HIBPCollector) GetStealerLogsByEmailDomain(ctx context.Context, domain string) (map[string][]string, error) {
	endpoint := fmt.Sprintf("/stealerlogsbyemaildomain/%s", url.QueryEscape(domain))
	
	var result map[string][]string
	err := c.makeRequest(ctx, endpoint, &result)
	return result, err
}

// Pwned Passwords 相關方法 (無需認證)
func (c *HIBPCollector) CheckPasswordHash(ctx context.Context, password string) (int, error) {
	// 計算 SHA-1 哈希
	hash := sha1.Sum([]byte(password))
	hashStr := strings.ToUpper(fmt.Sprintf("%x", hash))
	
	// 取前5位字符
	prefix := hashStr[:5]
	suffix := hashStr[5:]
	
	// 查詢 HIBP API
	endpoint := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", prefix)
	
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return 0, err
	}
	
	// 添加 padding 請求頭
	req.Header.Set("Add-Padding", "true")
	req.Header.Set("User-Agent", c.userAgent)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}
	
	// 解析響應
	body := make([]byte, 0)
	_, err = resp.Body.Read(body)
	if err != nil {
		return 0, err
	}
	
	// 查找匹配的哈希後綴
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 && strings.ToUpper(parts[0]) == suffix {
			count, _ := strconv.Atoi(parts[1])
			return count, nil
		}
	}
	
	return 0, nil
}

// 批量處理方法
func (c *HIBPCollector) ProcessAccountBreaches(ctx context.Context, account string) error {
	breaches, err := c.CheckAccountBreaches(ctx, account, true)
	if err != nil {
		return err
	}
	
	// 為每個泄露事件創建威脅情報
	for _, breach := range breaches {
		threatReq := &dto.ThreatIntelligenceCreateRequest{
			IPAddress:       "0.0.0.0", // 域名相關威脅，使用佔位符 IP
			Domain:          &breach.Domain,
			ThreatType:      "other",
			Severity:        c.determineBreachSeverity(breach),
			ConfidenceScore: 100, // HIBP 數據可信度很高
			Description:     &breach.Description,
			Source:          "HaveIBeenPwned",
			ExternalID:      &breach.Name,
			Tags:            []string{"data-breach", "hibp"},
			Metadata: map[string]interface{}{
				"breach_name":     breach.Name,
				"breach_date":     breach.BreachDate,
				"pwn_count":       breach.PwnCount,
				"data_classes":    breach.DataClasses,
				"is_verified":     breach.IsVerified,
				"is_sensitive":    breach.IsSensitive,
				"is_retired":      breach.IsRetired,
				"is_spam_list":    breach.IsSpamList,
				"is_malware":      breach.IsMalware,
				"is_stealer_log":  breach.IsStealerLog,
			},
		}
		
		_, err := c.service.CreateThreat(ctx, threatReq)
		if err != nil {
			// 記錄錯誤但繼續處理其他事件
			fmt.Printf("Error creating threat for breach %s: %v\n", breach.Name, err)
		}
	}
	
	return nil
}

// 輔助方法
func (c *HIBPCollector) makeRequest(ctx context.Context, endpoint string, result interface{}) error {
	url := c.baseURL + endpoint
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	
	// 設置必要的請求頭
	req.Header.Set("hibp-api-key", c.apiKey)
	req.Header.Set("user-agent", c.userAgent)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// 處理不同的響應狀態碼
	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(result)
	case http.StatusNotFound:
		return nil // 沒有找到數據，返回空結果
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized: invalid API key")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden: missing user agent")
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("retry-after")
		return fmt.Errorf("rate limit exceeded, retry after: %s", retryAfter)
	case http.StatusServiceUnavailable:
		return fmt.Errorf("service unavailable")
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *HIBPCollector) determineBreachSeverity(breach HIBPBreach) string {
	// 根據泄露事件的特性判斷嚴重程度
	if breach.IsSensitive {
		return "critical"
	}
	if breach.PwnCount > 10000000 { // 1000萬以上
		return "critical"
	}
	if breach.PwnCount > 1000000 { // 100萬以上
		return "high"
	}
	if breach.PwnCount > 100000 { // 10萬以上
		return "medium"
	}
	return "low"
}

// 健康檢查
func (c *HIBPCollector) HealthCheck(ctx context.Context) error {
	_, err := c.GetSubscriptionStatus(ctx)
	return err
}

// 獲取收集器統計信息
func (c *HIBPCollector) GetStats(ctx context.Context) (map[string]interface{}, error) {
	status, err := c.GetSubscriptionStatus(ctx)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"subscription_name":           status.Name,
		"pwn_count":                   status.PwnCount,
		"breach_count":                status.BreachCount,
		"paste_count":                 status.PasteCount,
		"domain_search_enabled":       status.DomainSearchEnabled,
		"domain_search_capacity":      status.DomainSearchCapacity,
		"domain_search_used":          status.DomainSearchUsed,
		"stealer_log_enabled":         status.StealerLogEnabled,
		"stealer_log_capacity":        status.StealerLogCapacity,
		"stealer_log_used":            status.StealerLogUsed,
		"next_renewal_date":           status.NextRenewalDate,
		"next_renewal_date_timezone":  status.NextRenewalDateTimezone,
	}, nil
}
