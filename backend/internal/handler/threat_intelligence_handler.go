package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
)

// ThreatIntelligenceHandler 威脅情報處理器
type ThreatIntelligenceHandler struct {
	service service.ThreatIntelligenceService
}

// NewThreatIntelligenceHandler 建立威脅情報處理器
func NewThreatIntelligenceHandler(service service.ThreatIntelligenceService) *ThreatIntelligenceHandler {
	return &ThreatIntelligenceHandler{service: service}
}

// CreateThreat 建立威脅情報
// @Summary 建立威脅情報
// @Description 建立新的威脅情報記錄
// @Tags Threat Intelligence
// @Accept json
// @Produce json
// @Param request body dto.ThreatIntelligenceCreateRequest true "威脅情報建立請求"
// @Success 201 {object} vo.BaseResponse{data=vo.ThreatIntelligenceVO} "建立成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats [post]
func (h *ThreatIntelligenceHandler) CreateThreat(c *gin.Context) {
	var req dto.ThreatIntelligenceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	threat, err := h.service.CreateThreat(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "CREATE_FAILED", "建立威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusCreated, "威脅情報建立成功", threat)
}

// GetThreat 取得威脅情報
// @Summary 取得威脅情報
// @Description 根據 ID 取得威脅情報詳情
// @Tags Threat Intelligence
// @Produce json
// @Param id path string true "威脅情報 ID" format(uuid)
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceVO} "取得成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 404 {object} vo.BaseResponse{error=vo.ErrorVO} "資源不存在"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/{id} [get]
func (h *ThreatIntelligenceHandler) GetThreat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_UUID", "無效的 UUID 格式", err)
		return
	}

	threat, err := h.service.GetThreatByID(c.Request.Context(), id)
	if err != nil {
		h.respondError(c, http.StatusNotFound, "NOT_FOUND", "威脅情報不存在", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "取得威脅情報成功", threat)
}

// UpdateThreat 更新威脅情報
// @Summary 更新威脅情報
// @Description 更新現有的威脅情報記錄
// @Tags Threat Intelligence
// @Accept json
// @Produce json
// @Param id path string true "威脅情報 ID" format(uuid)
// @Param request body dto.ThreatIntelligenceUpdateRequest true "威脅情報更新請求"
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceVO} "更新成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 404 {object} vo.BaseResponse{error=vo.ErrorVO} "資源不存在"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/{id} [put]
func (h *ThreatIntelligenceHandler) UpdateThreat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_UUID", "無效的 UUID 格式", err)
		return
	}

	var req dto.ThreatIntelligenceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	threat, err := h.service.UpdateThreat(c.Request.Context(), id, &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "UPDATE_FAILED", "更新威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "威脅情報更新成功", threat)
}

// DeleteThreat 刪除威脅情報
// @Summary 刪除威脅情報
// @Description 刪除指定的威脅情報記錄
// @Tags Threat Intelligence
// @Produce json
// @Param id path string true "威脅情報 ID" format(uuid)
// @Success 200 {object} vo.BaseResponse "刪除成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 404 {object} vo.BaseResponse{error=vo.ErrorVO} "資源不存在"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/{id} [delete]
func (h *ThreatIntelligenceHandler) DeleteThreat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_UUID", "無效的 UUID 格式", err)
		return
	}

	err = h.service.DeleteThreat(c.Request.Context(), id)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "DELETE_FAILED", "刪除威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "威脅情報刪除成功", nil)
}

// ListThreats 取得威脅情報列表
// @Summary 取得威脅情報列表
// @Description 取得威脅情報列表，支援篩選、分頁和排序
// @Tags Threat Intelligence
// @Produce json
// @Param ip_address query string false "IP 地址"
// @Param domain query string false "域名"
// @Param threat_type query string false "威脅類型" Enums(malware, phishing, spam, botnet, scanner, ddos, bruteforce, other)
// @Param severity query string false "嚴重程度" Enums(low, medium, high, critical)
// @Param source query string false "來源"
// @Param country_code query string false "國家代碼"
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁大小" default(20)
// @Param sort_by query string false "排序欄位" default(created_at) Enums(created_at, updated_at, last_seen, confidence_score)
// @Param sort_order query string false "排序順序" default(desc) Enums(asc, desc)
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceListVO} "取得成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats [get]
func (h *ThreatIntelligenceHandler) ListThreats(c *gin.Context) {
	var req dto.ThreatIntelligenceQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_QUERY", "查詢參數格式錯誤", err)
		return
	}

	threats, err := h.service.ListThreats(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "LIST_FAILED", "取得威脅情報列表失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "取得威脅情報列表成功", threats)
}

// LookupIP IP 查詢
// @Summary IP 威脅查詢
// @Description 查詢指定 IP 地址的威脅情報
// @Tags Threat Intelligence
// @Produce json
// @Param ip_address query string true "IP 地址"
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceIPLookupVO} "查詢成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/lookup/ip [get]
func (h *ThreatIntelligenceHandler) LookupIP(c *gin.Context) {
	var req dto.ThreatIntelligenceIPLookupRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_QUERY", "查詢參數格式錯誤", err)
		return
	}

	result, err := h.service.LookupIP(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "LOOKUP_FAILED", "IP 查詢失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "IP 查詢成功", result)
}

// LookupDomain 域名查詢
// @Summary 域名威脅查詢
// @Description 查詢指定域名的威脅情報
// @Tags Threat Intelligence
// @Produce json
// @Param domain query string true "域名"
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceDomainLookupVO} "查詢成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/lookup/domain [get]
func (h *ThreatIntelligenceHandler) LookupDomain(c *gin.Context) {
	var req dto.ThreatIntelligenceDomainLookupRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_QUERY", "查詢參數格式錯誤", err)
		return
	}

	result, err := h.service.LookupDomain(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "LOOKUP_FAILED", "域名查詢失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "域名查詢成功", result)
}

// GetStats 取得統計資料
// @Summary 取得威脅情報統計
// @Description 取得威脅情報的統計資料
// @Tags Threat Intelligence
// @Produce json
// @Param group_by query string false "分組依據" Enums(threat_type, severity, source, country_code)
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceStatsVO} "取得成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/stats [get]
func (h *ThreatIntelligenceHandler) GetStats(c *gin.Context) {
	var req dto.ThreatIntelligenceStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_QUERY", "查詢參數格式錯誤", err)
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "STATS_FAILED", "取得統計資料失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "取得統計資料成功", stats)
}

// BulkCreateThreats 批量建立威脅情報
// @Summary 批量建立威脅情報
// @Description 批量建立多個威脅情報記錄
// @Tags Threat Intelligence
// @Accept json
// @Produce json
// @Param request body dto.ThreatIntelligenceBulkCreateRequest true "批量建立請求"
// @Success 200 {object} vo.BaseResponse{data=vo.ThreatIntelligenceBulkCreateVO} "建立成功"
// @Failure 400 {object} vo.BaseResponse{error=vo.ErrorVO} "請求參數錯誤"
// @Failure 500 {object} vo.BaseResponse{error=vo.ErrorVO} "內部服務器錯誤"
// @Router /api/v1/threats/bulk [post]
func (h *ThreatIntelligenceHandler) BulkCreateThreats(c *gin.Context) {
	var req dto.ThreatIntelligenceBulkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	result, err := h.service.BulkCreateThreats(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "BULK_CREATE_FAILED", "批量建立威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "批量建立威脅情報完成", result)
}

// RegisterRoutes 註冊路由
func (h *ThreatIntelligenceHandler) RegisterRoutes(router *gin.RouterGroup) {
	threats := router.Group("/threats")
	{
		threats.POST("", h.CreateThreat)
		threats.GET("", h.ListThreats)
		threats.GET("/:id", h.GetThreat)
		threats.PUT("/:id", h.UpdateThreat)
		threats.DELETE("/:id", h.DeleteThreat)
		threats.POST("/bulk", h.BulkCreateThreats)
		
		// 查詢 API
		lookup := threats.Group("/lookup")
		{
			lookup.GET("/ip", h.LookupIP)
			lookup.GET("/domain", h.LookupDomain)
		}
		
		// 統計 API
		threats.GET("/stats", h.GetStats)
	}
}

// respondSuccess 回應成功
func (h *ThreatIntelligenceHandler) respondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := vo.BaseResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: h.getRequestID(c),
	}
	c.JSON(statusCode, response)
}

// respondError 回應錯誤
func (h *ThreatIntelligenceHandler) respondError(c *gin.Context, statusCode int, code string, message string, err error) {
	errorVO := &vo.ErrorVO{
		Code:    code,
		Message: message,
	}
	
	if err != nil {
		errorVO.Details = err.Error()
	}

	response := vo.BaseResponse{
		Success:   false,
		Message:   message,
		Error:     errorVO,
		Timestamp: time.Now(),
		RequestID: h.getRequestID(c),
	}
	c.JSON(statusCode, response)
}

// getRequestID 取得請求 ID
func (h *ThreatIntelligenceHandler) getRequestID(c *gin.Context) string {
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		return requestID
	}
	if requestID := c.GetString("request_id"); requestID != "" {
		return requestID
	}
	return ""
}

// parseIntParam 解析整數參數
func (h *ThreatIntelligenceHandler) parseIntParam(c *gin.Context, param string, defaultValue int) int {
	if valueStr := c.Query(param); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// SearchThreats 搜尋威脅情報
func (h *ThreatIntelligenceHandler) SearchThreats(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		h.respondError(c, http.StatusBadRequest, "INVALID_QUERY", "搜尋查詢不能為空", nil)
		return
	}

	page := h.parseIntParam(c, "page", 1)
	limit := h.parseIntParam(c, "limit", 10)

	threats, total, err := h.service.SearchThreats(c.Request.Context(), query, page, limit)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "SEARCH_FAILED", "搜尋威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "搜尋完成", gin.H{
		"threats": threats,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// GetStatistics 取得統計資訊
func (h *ThreatIntelligenceHandler) GetStatistics(c *gin.Context) {
	stats, err := h.service.GetStatistics(c.Request.Context())
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "STATS_FAILED", "取得統計資訊失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "統計資訊取得成功", stats)
}

// BulkUpdateThreats 批量更新威脅情報
func (h *ThreatIntelligenceHandler) BulkUpdateThreats(c *gin.Context) {
	var req dto.ThreatIntelligenceBulkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	result, err := h.service.BulkUpdateThreats(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "BULK_UPDATE_FAILED", "批量更新威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "批量更新威脅情報完成", result)
}

// BulkDeleteThreats 批量刪除威脅情報
func (h *ThreatIntelligenceHandler) BulkDeleteThreats(c *gin.Context) {
	var req dto.ThreatIntelligenceBulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "請求參數格式錯誤", err)
		return
	}

	result, err := h.service.BulkDeleteThreats(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, "BULK_DELETE_FAILED", "批量刪除威脅情報失敗", err)
		return
	}

	h.respondSuccess(c, http.StatusOK, "批量刪除威脅情報完成", result)
} 