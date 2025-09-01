package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
)

// QuantumJobHandler 處理量子任務相關請求
type QuantumJobHandler struct {
	quantumJobService *service.QuantumJobService
}

// NewQuantumJobHandler 創建量子任務處理器
func NewQuantumJobHandler(quantumJobService *service.QuantumJobService) *QuantumJobHandler {
	return &QuantumJobHandler{
		quantumJobService: quantumJobService,
	}
}

// RegisterRoutes 註冊路由
func (h *QuantumJobHandler) RegisterRoutes(router *gin.RouterGroup) {
	quantumGroup := router.Group("/quantum-jobs")
	{
		quantumGroup.POST("", h.SubmitJob)
		quantumGroup.GET("", h.ListJobs)
		quantumGroup.GET("/:job_id", h.GetJob)
		quantumGroup.PUT("/:job_id", h.UpdateJob)
		quantumGroup.DELETE("/:job_id", h.CancelJob)
		quantumGroup.POST("/batch", h.BatchSubmitJobs)
		quantumGroup.POST("/threat-analysis", h.SubmitThreatAnalysis)
	}

	// 內部API，需要額外的認證
	internalGroup := router.Group("/internal/quantum-jobs")
	{
		internalGroup.POST("/next", h.GetNextPendingJob)
		internalGroup.POST("/complete", h.CompleteJob)
		internalGroup.POST("/fail", h.FailJob)
	}
}

// SubmitJob 提交新的量子任務
// @Summary 提交新的量子任務
// @Description 提交新的量子計算任務進行排程執行
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param request body dto.SubmitQuantumJobRequest true "任務資訊"
// @Success 201 {object} vo.SubmitQuantumJobResponse "任務已提交"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Router /api/v1/quantum-jobs [post]
func (h *QuantumJobHandler) SubmitJob(c *gin.Context) {
	var req dto.SubmitQuantumJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能提交量子任務",
		})
		return
	}

	// 調用服務層
	jobID, err := h.quantumJobService.SubmitJob(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Error:   "Failed to submit job",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusCreated, vo.SubmitQuantumJobResponse{
		JobID:     jobID,
		Status:    "pending",
		Message:   "量子任務已成功提交，等待執行",
		CreatedAt: h.quantumJobService.GetJobCreationTime(c, jobID),
	})
}

// GetJob 獲取特定量子任務的詳細信息
// @Summary 獲取量子任務詳情
// @Description 根據任務ID獲取量子任務的詳細信息
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param job_id path string true "任務ID"
// @Success 200 {object} vo.QuantumJobDetailResponse "任務詳情"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Failure 404 {object} vo.ErrorResponse "任務不存在"
// @Router /api/v1/quantum-jobs/{job_id} [get]
func (h *QuantumJobHandler) GetJob(c *gin.Context) {
	var req dto.GetQuantumJobRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能獲取任務詳情",
		})
		return
	}

	// 調用服務層
	jobDetail, err := h.quantumJobService.GetJob(c, userID, req.JobID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "record not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, vo.ErrorResponse{
			Error:   "Failed to get job",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, jobDetail)
}

// ListJobs 獲取量子任務列表
// @Summary 列出量子任務
// @Description 根據過濾條件列出用戶的量子任務
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param status query string false "狀態過濾" Enums(pending, running, completed, failed)
// @Param page query int false "頁碼" default(1) minimum(1)
// @Param page_size query int false "每頁大小" default(20) minimum(1) maximum(100)
// @Param sort_by query string false "排序欄位" Enums(created_at, priority, title, status) default(created_at)
// @Param sort_order query string false "排序方式" Enums(asc, desc) default(desc)
// @Success 200 {object} vo.ListQuantumJobsResponse "任務列表"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Router /api/v1/quantum-jobs [get]
func (h *QuantumJobHandler) ListJobs(c *gin.Context) {
	var req dto.ListQuantumJobsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 設置默認值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能獲取任務列表",
		})
		return
	}

	// 調用服務層
	jobs, pagination, err := h.quantumJobService.ListJobs(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Error:   "Failed to list jobs",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, vo.ListQuantumJobsResponse{
		Jobs:       jobs,
		Pagination: *pagination,
	})
}

// UpdateJob 更新量子任務
// @Summary 更新量子任務
// @Description 更新量子任務的標題、描述、優先級等資訊
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param job_id path string true "任務ID"
// @Param request body dto.UpdateQuantumJobRequest true "更新資訊"
// @Success 200 {object} vo.QuantumJobResponse "更新後的任務"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Failure 404 {object} vo.ErrorResponse "任務不存在"
// @Failure 409 {object} vo.ErrorResponse "任務已開始執行，無法更新"
// @Router /api/v1/quantum-jobs/{job_id} [put]
func (h *QuantumJobHandler) UpdateJob(c *gin.Context) {
	var req dto.UpdateQuantumJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}
	
	// 綁定URI參數
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid URI parameters",
			Message: err.Error(),
		})
		return
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能更新任務",
		})
		return
	}

	// 調用服務層
	jobResponse, err := h.quantumJobService.UpdateJob(c, userID, req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "record not found":
			status = http.StatusNotFound
		case "cannot update job that is already running or completed":
			status = http.StatusConflict
		}
		c.JSON(status, vo.ErrorResponse{
			Error:   "Failed to update job",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, jobResponse)
}

// CancelJob 取消量子任務
// @Summary 取消量子任務
// @Description 取消一個尚未執行的量子任務
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param job_id path string true "任務ID"
// @Param request body dto.CancelQuantumJobRequest true "取消資訊"
// @Success 200 {object} vo.SuccessResponse "取消成功"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Failure 404 {object} vo.ErrorResponse "任務不存在"
// @Failure 409 {object} vo.ErrorResponse "任務已開始執行，無法取消"
// @Router /api/v1/quantum-jobs/{job_id} [delete]
func (h *QuantumJobHandler) CancelJob(c *gin.Context) {
	var req dto.CancelQuantumJobRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid URI parameters",
			Message: err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON綁定失敗不是錯誤，因為reason是可選的
		// 繼續處理
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能取消任務",
		})
		return
	}

	// 調用服務層
	err := h.quantumJobService.CancelJob(c, userID, req.JobID, req.Reason)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "record not found":
			status = http.StatusNotFound
		case "cannot cancel job that is already running or completed":
			status = http.StatusConflict
		}
		c.JSON(status, vo.ErrorResponse{
			Error:   "Failed to cancel job",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, vo.SuccessResponse{
		Success: true,
		Message: "量子任務已成功取消",
	})
}

// BatchSubmitJobs 批次提交量子任務
// @Summary 批次提交量子任務
// @Description 一次提交多個量子任務
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param request body dto.BatchSubmitQuantumJobsRequest true "批次任務資訊"
// @Success 201 {object} vo.BatchSubmitQuantumJobsResponse "批次提交結果"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Router /api/v1/quantum-jobs/batch [post]
func (h *QuantumJobHandler) BatchSubmitJobs(c *gin.Context) {
	var req dto.BatchSubmitQuantumJobsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能批次提交任務",
		})
		return
	}

	// 調用服務層
	response, err := h.quantumJobService.BatchSubmitJobs(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Error:   "Failed to submit batch jobs",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusCreated, response)
}

// SubmitThreatAnalysis 提交威脅分析任務
// @Summary 提交威脅分析任務
// @Description 提交特定格式的威脅分析量子任務
// @Tags 量子任務
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param request body dto.ThreatAnalysisRequest true "威脅分析請求"
// @Success 201 {object} vo.ThreatAnalysisSubmitResponse "任務已提交"
// @Failure 400 {object} vo.ErrorResponse "參數錯誤"
// @Failure 401 {object} vo.ErrorResponse "未授權"
// @Router /api/v1/quantum-jobs/threat-analysis [post]
func (h *QuantumJobHandler) SubmitThreatAnalysis(c *gin.Context) {
	var req dto.ThreatAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 獲取當前用戶ID
	userID, exists := getUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "您需要登入才能提交威脅分析",
		})
		return
	}

	// 調用服務層
	response, err := h.quantumJobService.SubmitThreatAnalysis(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Error:   "Failed to submit threat analysis",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusCreated, response)
}

// ---- 以下為內部API，需要額外的認證 ----

// GetNextPendingJob 獲取下一個待處理的任務（內部API）
func (h *QuantumJobHandler) GetNextPendingJob(c *gin.Context) {
	var req dto.GetNextPendingJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 驗證內部API密鑰
	if !h.validateInternalAPIKey(c, req.Secret) {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "無效的API密鑰",
		})
		return
	}

	// 調用服務層
	jobResponse, err := h.quantumJobService.GetNextPendingJob(c, req.WorkerID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "no pending jobs available" {
			status = http.StatusNotFound
		}
		c.JSON(status, vo.ErrorResponse{
			Error:   "Failed to get next pending job",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, jobResponse)
}

// CompleteJob 標記任務為完成（內部API）
func (h *QuantumJobHandler) CompleteJob(c *gin.Context) {
	var req dto.CompleteQuantumJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 驗證內部API密鑰
	if !h.validateInternalAPIKey(c, c.GetHeader("X-API-KEY")) {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "無效的API密鑰",
		})
		return
	}

	// 調用服務層
	err := h.quantumJobService.CompleteJob(c, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "record not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, vo.ErrorResponse{
			Error:   "Failed to complete job",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, vo.SuccessResponse{
		Success: true,
		Message: "量子任務已標記為完成",
	})
}

// FailJob 標記任務為失敗（內部API）
func (h *QuantumJobHandler) FailJob(c *gin.Context) {
	var req dto.FailQuantumJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// 驗證內部API密鑰
	if !h.validateInternalAPIKey(c, c.GetHeader("X-API-KEY")) {
		c.JSON(http.StatusUnauthorized, vo.ErrorResponse{
			Error:   "Unauthorized",
			Message: "無效的API密鑰",
		})
		return
	}

	// 調用服務層
	err := h.quantumJobService.FailJob(c, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "record not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, vo.ErrorResponse{
			Error:   "Failed to mark job as failed",
			Message: err.Error(),
		})
		return
	}

	// 返回成功響應
	c.JSON(http.StatusOK, vo.SuccessResponse{
		Success: true,
		Message: "量子任務已標記為失敗",
	})
}

// 輔助函數：從上下文獲取用戶ID
func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	
	id, ok := userID.(uuid.UUID)
	if !ok {
		idStr, ok := userID.(string)
		if !ok {
			return uuid.Nil, false
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, false
		}
		return id, true
	}
	
	return id, true
}

// 驗證內部API密鑰
func (h *QuantumJobHandler) validateInternalAPIKey(c *gin.Context, key string) bool {
	// TODO: 實現內部API密鑰驗證
	// 這裡應該檢查環境變數或配置文件中的密鑰
	expectedKey := "your-internal-api-key" // 實際應用中從配置中獲取
	return key == expectedKey
}
