package handler

import (
	"net/http"
	"time"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/collector"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// HIBPHandler 處理 HIBP API 相關請求
type HIBPHandler struct {
	hibpCollector *collector.HIBPCollector
	threatService service.ThreatIntelligenceService
}

// NewHIBPHandler 創建 HIBP 處理器
func NewHIBPHandler(hibpCollector *collector.HIBPCollector, threatService service.ThreatIntelligenceService) *HIBPHandler {
	return &HIBPHandler{
		hibpCollector: hibpCollector,
		threatService: threatService,
	}
}

// CheckAccountBreaches 檢查帳戶泄露
func (h *HIBPHandler) CheckAccountBreaches(c *gin.Context) {
	account := c.Param("account")
	if account == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_ACCOUNT", "Account parameter is required", nil)
		return
	}

	includeUnverified := c.Query("include_unverified") == "true"

	breaches, err := h.hibpCollector.CheckAccountBreaches(c.Request.Context(), account, includeUnverified)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to check account breaches", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Account breaches retrieved successfully", gin.H{
		"account":  account,
		"breaches": breaches,
		"count":    len(breaches),
	})
}

// GetDomainBreaches 獲取域名泄露
func (h *HIBPHandler) GetDomainBreaches(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_DOMAIN", "Domain parameter is required", nil)
		return
	}

	domainBreaches, err := h.hibpCollector.GetDomainBreaches(c.Request.Context(), domain)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get domain breaches", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Domain breaches retrieved successfully", gin.H{
		"domain":   domain,
		"breaches": domainBreaches,
		"count":    len(domainBreaches),
	})
}

// GetSubscribedDomains 獲取已訂閱域名
func (h *HIBPHandler) GetSubscribedDomains(c *gin.Context) {
	domains, err := h.hibpCollector.GetSubscribedDomains(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get subscribed domains", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Subscribed domains retrieved successfully", gin.H{
		"domains": domains,
		"count":   len(domains),
	})
}

// GetAllBreaches 獲取所有泄露事件
func (h *HIBPHandler) GetAllBreaches(c *gin.Context) {
	domain := c.Query("domain")

	breaches, err := h.hibpCollector.GetAllBreaches(c.Request.Context(), domain)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get all breaches", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "All breaches retrieved successfully", gin.H{
		"breaches": breaches,
		"count":    len(breaches),
		"domain":   domain,
	})
}

// GetBreachByName 根據名稱獲取泄露事件
func (h *HIBPHandler) GetBreachByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_BREACH_NAME", "Breach name parameter is required", nil)
		return
	}

	breach, err := h.hibpCollector.GetBreachByName(c.Request.Context(), name)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get breach by name", err)
		return
	}

	if breach == nil {
		h.respondError(c, http.StatusNotFound, "BREACH_NOT_FOUND", "Breach not found", nil)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Breach retrieved successfully", breach)
}

// GetLatestBreach 獲取最新泄露事件
func (h *HIBPHandler) GetLatestBreach(c *gin.Context) {
	breach, err := h.hibpCollector.GetLatestBreach(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get latest breach", err)
		return
	}

	if breach == nil {
		h.respondError(c, http.StatusNotFound, "NO_BREACHES_FOUND", "No breaches found", nil)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Latest breach retrieved successfully", breach)
}

// GetDataClasses 獲取數據類別
func (h *HIBPHandler) GetDataClasses(c *gin.Context) {
	dataClasses, err := h.hibpCollector.GetDataClasses(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get data classes", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Data classes retrieved successfully", gin.H{
		"data_classes": dataClasses,
		"count":        len(dataClasses),
	})
}

// GetAccountPastes 獲取帳戶 Paste 記錄
func (h *HIBPHandler) GetAccountPastes(c *gin.Context) {
	account := c.Param("account")
	if account == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_ACCOUNT", "Account parameter is required", nil)
		return
	}

	pastes, err := h.hibpCollector.GetAccountPastes(c.Request.Context(), account)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get account pastes", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Account pastes retrieved successfully", gin.H{
		"account": account,
		"pastes":  pastes,
		"count":   len(pastes),
	})
}

// GetSubscriptionStatus 獲取訂閱狀態
func (h *HIBPHandler) GetSubscriptionStatus(c *gin.Context) {
	status, err := h.hibpCollector.GetSubscriptionStatus(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get subscription status", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Subscription status retrieved successfully", status)
}

// CheckPasswordHash 檢查密碼哈希
func (h *HIBPHandler) CheckPasswordHash(c *gin.Context) {
	password := c.Query("password")
	if password == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_PASSWORD", "Password parameter is required", nil)
		return
	}

	count, err := h.hibpCollector.CheckPasswordHash(c.Request.Context(), password)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to check password hash", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Password hash check completed", gin.H{
		"pwned_count": count,
		"is_pwned":    count > 0,
	})
}

// Stealer Logs 相關端點
func (h *HIBPHandler) GetStealerLogsByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_EMAIL", "Email parameter is required", nil)
		return
	}

	domains, err := h.hibpCollector.GetStealerLogsByEmail(c.Request.Context(), email)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get stealer logs by email", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Stealer logs by email retrieved successfully", gin.H{
		"email":   email,
		"domains": domains,
		"count":   len(domains),
	})
}

func (h *HIBPHandler) GetStealerLogsByWebsiteDomain(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_DOMAIN", "Domain parameter is required", nil)
		return
	}

	emails, err := h.hibpCollector.GetStealerLogsByWebsiteDomain(c.Request.Context(), domain)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get stealer logs by website domain", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Stealer logs by website domain retrieved successfully", gin.H{
		"domain": domain,
		"emails": emails,
		"count":  len(emails),
	})
}

func (h *HIBPHandler) GetStealerLogsByEmailDomain(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_DOMAIN", "Domain parameter is required", nil)
		return
	}

	result, err := h.hibpCollector.GetStealerLogsByEmailDomain(c.Request.Context(), domain)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get stealer logs by email domain", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Stealer logs by email domain retrieved successfully", gin.H{
		"domain": domain,
		"result": result,
		"count":  len(result),
	})
}

// ProcessAccountBreaches 處理帳戶泄露並創建威脅情報
func (h *HIBPHandler) ProcessAccountBreaches(c *gin.Context) {
	account := c.Param("account")
	if account == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_ACCOUNT", "Account parameter is required", nil)
		return
	}

	err := h.hibpCollector.ProcessAccountBreaches(c.Request.Context(), account)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_PROCESSING_ERROR", "Failed to process account breaches", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "Account breaches processed successfully", gin.H{
		"account": account,
		"status":  "processed",
	})
}

// GetHIBPStats 獲取 HIBP 統計信息
func (h *HIBPHandler) GetHIBPStats(c *gin.Context) {
	stats, err := h.hibpCollector.GetStats(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "HIBP_API_ERROR", "Failed to get HIBP stats", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "HIBP stats retrieved successfully", stats)
}

// HealthCheck HIBP 健康檢查
func (h *HIBPHandler) HealthCheck(c *gin.Context) {
	err := h.hibpCollector.HealthCheck(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusServiceUnavailable, "HIBP_UNHEALTHY", "HIBP service is unhealthy", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "HIBP service is healthy", gin.H{
		"status": "healthy",
	})
}

// RegisterRoutes 註冊路由
func (h *HIBPHandler) RegisterRoutes(router *gin.RouterGroup) {
	hibp := router.Group("/hibp")
	{
		// 帳戶相關
		hibp.GET("/account/:account/breaches", h.CheckAccountBreaches)
		hibp.POST("/account/:account/process", h.ProcessAccountBreaches)
		hibp.GET("/account/:account/pastes", h.GetAccountPastes)

		// 域名相關
		hibp.GET("/domain/:domain/breaches", h.GetDomainBreaches)
		hibp.GET("/domains/subscribed", h.GetSubscribedDomains)

		// 泄露事件相關
		hibp.GET("/breaches", h.GetAllBreaches)
		hibp.GET("/breach/:name", h.GetBreachByName)
		hibp.GET("/breach/latest", h.GetLatestBreach)
		hibp.GET("/dataclasses", h.GetDataClasses)

		// Stealer Logs (需要高級訂閱)
		hibp.GET("/stealer/email/:email", h.GetStealerLogsByEmail)
		hibp.GET("/stealer/website/:domain", h.GetStealerLogsByWebsiteDomain)
		hibp.GET("/stealer/emaildomain/:domain", h.GetStealerLogsByEmailDomain)

		// 密碼檢查 (無需認證)
		hibp.GET("/password/check", h.CheckPasswordHash)

		// 系統相關
		hibp.GET("/subscription/status", h.GetSubscriptionStatus)
		hibp.GET("/stats", h.GetHIBPStats)
		hibp.GET("/health", h.HealthCheck)
	}
}

// 輔助方法
func (h *HIBPHandler) respondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success":   true,
		"message":   message,
		"data":      data,
		"timestamp": time.Now().UTC(),
		"request_id": c.GetString("request_id"),
	})
}

func (h *HIBPHandler) respondError(c *gin.Context, statusCode int, code string, message string, err error) {
	errorResponse := gin.H{
		"code":    code,
		"message": message,
	}

	if err != nil {
		errorResponse["details"] = err.Error()
	}

	c.JSON(statusCode, gin.H{
		"success":   false,
		"message":   message,
		"error":     errorResponse,
		"timestamp": time.Now().UTC(),
		"request_id": c.GetString("request_id"),
	})
}
