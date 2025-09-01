package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/collector"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
)

// CollectorHandler 收集器處理器
type CollectorHandler struct {
	threatIntelService service.ThreatIntelligenceService
}

// NewCollectorHandler 建立收集器處理器
func NewCollectorHandler(threatIntelService service.ThreatIntelligenceService) *CollectorHandler {
	return &CollectorHandler{
		threatIntelService: threatIntelService,
	}
}

// CollectIPRequest IP 收集請求
type CollectIPRequest struct {
	IPAddress string `json:"ip_address" binding:"required,ip"`
}

// CollectBulkIPRequest 批量 IP 收集請求
type CollectBulkIPRequest struct {
	IPAddresses []string `json:"ip_addresses" binding:"required,min=1,max=50"`
}

// CollectIPThreatIntel 收集單一 IP 威脅情報
// @Summary 收集 IP 威脅情報
// @Description 從 AbuseIPDB 收集指定 IP 的威脅情報
// @Tags Collector
// @Accept json
// @Produce json
// @Param request body CollectIPRequest true "IP 收集請求"
// @Success 200 {object} vo.BaseResponse "收集成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "收集失敗"
// @Router /api/v1/collector/ip [post]
func (h *CollectorHandler) CollectIPThreatIntel(c *gin.Context) {
	var req CollectIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	// 檢查 AbuseIPDB API 金鑰
	apiKey := os.Getenv("ABUSEIPDB_API_KEY")
	if apiKey == "" {
		h.respondError(c, http.StatusInternalServerError, "CONFIG_ERROR", "AbuseIPDB API 金鑰未設定", nil)
		return
	}

	// 建立收集器
	abuseCollector := collector.NewAbuseIPDBCollector(apiKey, h.threatIntelService)

	// 執行收集
	err := abuseCollector.CollectIPThreatIntel(c.Request.Context(), req.IPAddress)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "COLLECT_FAILED", "威脅情報收集失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "威脅情報收集成功", gin.H{
		"ip_address": req.IPAddress,
		"status":     "collected",
	})
}

// CollectBulkIPThreatIntel 批量收集 IP 威脅情報
// @Summary 批量收集 IP 威脅情報
// @Description 從 AbuseIPDB 批量收集多個 IP 的威脅情報
// @Tags Collector
// @Accept json
// @Produce json
// @Param request body CollectBulkIPRequest true "批量 IP 收集請求"
// @Success 200 {object} vo.BaseResponse "收集完成"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "收集失敗"
// @Router /api/v1/collector/bulk-ip [post]
func (h *CollectorHandler) CollectBulkIPThreatIntel(c *gin.Context) {
	var req CollectBulkIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	// 檢查 AbuseIPDB API 金鑰
	apiKey := os.Getenv("ABUSEIPDB_API_KEY")
	if apiKey == "" {
		h.respondError(c, http.StatusInternalServerError, "CONFIG_ERROR", "AbuseIPDB API 金鑰未設定", nil)
		return
	}

	// 建立收集器
	abuseCollector := collector.NewAbuseIPDBCollector(apiKey, h.threatIntelService)

	// 執行批量收集
	successful, failed := abuseCollector.CollectBulkIPThreatIntel(c.Request.Context(), req.IPAddresses)

	h.respondSuccess(c, http.StatusOK, "批量威脅情報收集完成", gin.H{
		"total":      len(req.IPAddresses),
		"successful": len(successful),
		"failed":     len(failed),
		"successful_ips": successful,
		"failed_ips":     failed,
	})
}

// respondSuccess 回傳成功回應
func (h *CollectorHandler) respondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := vo.BaseResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: h.getRequestID(c),
	}
	c.JSON(statusCode, response)
}

// respondError 回傳錯誤回應
func (h *CollectorHandler) respondError(c *gin.Context, statusCode int, code string, message string, err error) {
	errorVO := vo.ErrorVO{
		Code:    code,
		Message: message,
	}
	
	if err != nil {
		errorVO.Details = err.Error()
	}

	response := vo.BaseResponse{
		Success:   false,
		Message:   "請求處理失敗",
		Error:     &errorVO,
		RequestID: h.getRequestID(c),
	}
	
	c.JSON(statusCode, response)
}

// getRequestID 取得請求 ID
func (h *CollectorHandler) getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
} 