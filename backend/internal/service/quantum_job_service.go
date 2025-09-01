package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/model"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// QuantumJobService 量子任務服務
type QuantumJobService struct {
	db *gorm.DB
}

// NewQuantumJobService 創建量子任務服務
func NewQuantumJobService(db *gorm.DB) *QuantumJobService {
	return &QuantumJobService{
		db: db,
	}
}

// SubmitJob 提交量子任務
func (s *QuantumJobService) SubmitJob(ctx context.Context, userID uuid.UUID, req dto.SubmitQuantumJobRequest) (uuid.UUID, error) {
	// 創建量子任務
	job := model.NewQuantumJob(userID, req.Title, req.InputParams)
	
	// 設置額外資訊
	if req.Description != "" {
		job.Description = req.Description
	}
	if req.Priority > 0 {
		job.Priority = req.Priority
	}
	if len(req.Tags) > 0 {
		job.Tags = req.Tags
	}
	if req.Notes != "" {
		job.Notes = req.Notes
	}
	if req.Source != "" {
		job.Source = req.Source
	}

	// 保存到資料庫
	if err := s.db.Create(job).Error; err != nil {
		logger.Error("Failed to create quantum job", logger.Fields{
			"error":    err.Error(),
			"user_id":  userID,
			"title":    req.Title,
		})
		return uuid.Nil, fmt.Errorf("failed to create quantum job: %w", err)
	}

	logger.Info("Quantum job submitted successfully", logger.Fields{
		"job_id":   job.ID,
		"user_id":  userID,
		"title":    req.Title,
		"priority": req.Priority,
	})

	return job.ID, nil
}

// GetJob 獲取量子任務詳情
func (s *QuantumJobService) GetJob(ctx context.Context, userID uuid.UUID, jobID uuid.UUID) (*vo.QuantumJobDetailResponse, error) {
	var job model.QuantumJob
	
	// 查詢任務，包含日誌
	err := s.db.Preload("Logs").Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("record not found")
		}
		logger.Error("Failed to get quantum job", logger.Fields{
			"error":   err.Error(),
			"job_id":  jobID,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get quantum job: %w", err)
	}

	// 轉換為 VO
	response := s.convertToDetailResponse(&job)
	return response, nil
}

// ListJobs 列出量子任務
func (s *QuantumJobService) ListJobs(ctx context.Context, userID uuid.UUID, req dto.ListQuantumJobsRequest) ([]vo.QuantumJobResponse, *vo.PaginationInfo, error) {
	var jobs []model.QuantumJob
	var total int64

	// 構建查詢
	query := s.db.Where("user_id = ?", userID)
	
	// 狀態過濾
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 計算總數
	if err := query.Model(&model.QuantumJob{}).Count(&total).Error; err != nil {
		logger.Error("Failed to count quantum jobs", logger.Fields{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, nil, fmt.Errorf("failed to count quantum jobs: %w", err)
	}

	// 排序
	orderBy := "created_at DESC"
	if req.SortBy != "" {
		orderBy = req.SortBy
		if req.SortOrder == "asc" {
			orderBy += " ASC"
		} else {
			orderBy += " DESC"
		}
	}

	// 分頁查詢
	offset := (req.Page - 1) * req.PageSize
	err := query.Order(orderBy).Offset(offset).Limit(req.PageSize).Find(&jobs).Error
	if err != nil {
		logger.Error("Failed to list quantum jobs", logger.Fields{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, nil, fmt.Errorf("failed to list quantum jobs: %w", err)
	}

	// 轉換為 VO
	responses := make([]vo.QuantumJobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = s.convertToResponse(&job)
	}

	// 計算分頁資訊
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	pagination := &vo.PaginationInfo{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return responses, pagination, nil
}

// UpdateJob 更新量子任務
func (s *QuantumJobService) UpdateJob(ctx context.Context, userID uuid.UUID, req dto.UpdateQuantumJobRequest) (*vo.QuantumJobResponse, error) {
	var job model.QuantumJob
	
	// 查詢任務
	err := s.db.Where("id = ? AND user_id = ?", req.JobID, userID).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to get quantum job: %w", err)
	}

	// 檢查是否可以更新
	if job.Status == model.JobStatusRunning || job.Status == model.JobStatusCompleted || job.Status == model.JobStatusFailed {
		return nil, fmt.Errorf("cannot update job that is already running or completed")
	}

	// 更新欄位
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if len(req.Tags) > 0 {
		updates["tags"] = req.Tags
	}
	if req.InputParams != nil {
		updates["input_params"] = req.InputParams
	}

	// 執行更新
	if err := s.db.Model(&job).Updates(updates).Error; err != nil {
		logger.Error("Failed to update quantum job", logger.Fields{
			"error":   err.Error(),
			"job_id":  req.JobID,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to update quantum job: %w", err)
	}

	// 重新查詢以獲取更新後的資料
	if err := s.db.First(&job, req.JobID).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated quantum job: %w", err)
	}

	response := s.convertToResponse(&job)
	return &response, nil
}

// CancelJob 取消量子任務
func (s *QuantumJobService) CancelJob(ctx context.Context, userID uuid.UUID, jobID uuid.UUID, reason string) error {
	var job model.QuantumJob
	
	// 查詢任務
	err := s.db.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("record not found")
		}
		return fmt.Errorf("failed to get quantum job: %w", err)
	}

	// 檢查是否可以取消
	if job.Status == model.JobStatusRunning || job.Status == model.JobStatusCompleted || job.Status == model.JobStatusFailed {
		return fmt.Errorf("cannot cancel job that is already running or completed")
	}

	// 更新狀態
	updates := map[string]interface{}{
		"status":        model.JobStatusFailed,
		"error_message": "Task cancelled by user",
	}
	if reason != "" {
		updates["error_message"] = fmt.Sprintf("Task cancelled by user: %s", reason)
	}

	if err := s.db.Model(&job).Updates(updates).Error; err != nil {
		logger.Error("Failed to cancel quantum job", logger.Fields{
			"error":   err.Error(),
			"job_id":  jobID,
			"user_id": userID,
		})
		return fmt.Errorf("failed to cancel quantum job: %w", err)
	}

	logger.Info("Quantum job cancelled successfully", logger.Fields{
		"job_id":  jobID,
		"user_id": userID,
		"reason":  reason,
	})

	return nil
}

// BatchSubmitJobs 批次提交量子任務
func (s *QuantumJobService) BatchSubmitJobs(ctx context.Context, userID uuid.UUID, req dto.BatchSubmitQuantumJobsRequest) (*vo.BatchSubmitQuantumJobsResponse, error) {
	response := &vo.BatchSubmitQuantumJobsResponse{
		JobIDs:   make([]uuid.UUID, 0),
		Failures: make([]vo.BatchSubmissionFailure, 0),
	}

	for i, jobReq := range req.Jobs {
		jobID, err := s.SubmitJob(ctx, userID, jobReq)
		if err != nil {
			response.FailedCount++
			response.Failures = append(response.Failures, vo.BatchSubmissionFailure{
				Index:   i,
				Message: err.Error(),
			})
		} else {
			response.SuccessCount++
			response.JobIDs = append(response.JobIDs, jobID)
		}
	}

	return response, nil
}

// SubmitThreatAnalysis 提交威脅分析任務
func (s *QuantumJobService) SubmitThreatAnalysis(ctx context.Context, userID uuid.UUID, req dto.ThreatAnalysisRequest) (*vo.ThreatAnalysisSubmitResponse, error) {
	// 構建輸入參數
	inputParams := map[string]interface{}{
		"data_sources":  req.DataSources,
		"threat_type":   req.ThreatType,
		"time_window":   req.TimeWindow,
		"use_simulator": req.UseSimulator,
	}

	// 創建任務請求
	jobReq := dto.SubmitQuantumJobRequest{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		InputParams: inputParams,
		Source:      "threat_analysis",
	}

	// 提交任務
	jobID, err := s.SubmitJob(ctx, userID, jobReq)
	if err != nil {
		return nil, err
	}

	// 估算執行時間
	estimatedTime := 300 // 預設5分鐘
	if !req.UseSimulator {
		estimatedTime = 600 // 真實設備需要更長時間
	}

	response := &vo.ThreatAnalysisSubmitResponse{
		JobID:           jobID,
		Status:          "pending",
		Title:           req.Title,
		Message:         "威脅分析任務已提交",
		EstimatedTime:   estimatedTime,
	}

	return response, nil
}

// GetNextPendingJob 獲取下一個待處理的任務（內部使用）
func (s *QuantumJobService) GetNextPendingJob(ctx context.Context, workerID string) (*vo.GetNextPendingJobResponse, error) {
	var job model.QuantumJob
	
	// 使用資料庫函數獲取下一個待處理任務
	err := s.db.Raw("SELECT * FROM get_next_pending_quantum_job()").Scan(&job).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no pending jobs available")
		}
		logger.Error("Failed to get next pending quantum job", logger.Fields{
			"error":     err.Error(),
			"worker_id": workerID,
		})
		return nil, fmt.Errorf("failed to get next pending quantum job: %w", err)
	}

	response := &vo.GetNextPendingJobResponse{
		JobID:       job.ID,
		InputParams: job.InputParams,
		Title:       job.Title,
		Priority:    job.Priority,
	}

	return response, nil
}

// CompleteJob 標記任務為完成（內部使用）
func (s *QuantumJobService) CompleteJob(ctx context.Context, req dto.CompleteQuantumJobRequest) error {
	var job model.QuantumJob
	
	// 查詢任務
	err := s.db.First(&job, req.JobID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("record not found")
		}
		return fmt.Errorf("failed to get quantum job: %w", err)
	}

	// 標記為完成
	job.MarkAsCompleted(req.Results, req.ConfidenceScore, req.IsMalicious)
	job.QuantumBackend = req.QuantumBackend

	if err := s.db.Save(&job).Error; err != nil {
		logger.Error("Failed to complete quantum job", logger.Fields{
			"error":   err.Error(),
			"job_id":  req.JobID,
		})
		return fmt.Errorf("failed to complete quantum job: %w", err)
	}

	logger.Info("Quantum job completed successfully", logger.Fields{
		"job_id":           req.JobID,
		"confidence_score": req.ConfidenceScore,
		"is_malicious":     req.IsMalicious,
		"execution_time":   req.ExecutionTime,
	})

	return nil
}

// FailJob 標記任務為失敗（內部使用）
func (s *QuantumJobService) FailJob(ctx context.Context, req dto.FailQuantumJobRequest) error {
	var job model.QuantumJob
	
	// 查詢任務
	err := s.db.First(&job, req.JobID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("record not found")
		}
		return fmt.Errorf("failed to get quantum job: %w", err)
	}

	// 標記為失敗
	job.MarkAsFailed(req.ErrorMessage)

	if err := s.db.Save(&job).Error; err != nil {
		logger.Error("Failed to mark quantum job as failed", logger.Fields{
			"error":   err.Error(),
			"job_id":  req.JobID,
		})
		return fmt.Errorf("failed to mark quantum job as failed: %w", err)
	}

	logger.Info("Quantum job marked as failed", logger.Fields{
		"job_id":        req.JobID,
		"error_message": req.ErrorMessage,
		"execution_time": req.ExecutionTime,
	})

	return nil
}

// GetJobCreationTime 獲取任務創建時間
func (s *QuantumJobService) GetJobCreationTime(ctx context.Context, jobID uuid.UUID) time.Time {
	var job model.QuantumJob
	if err := s.db.Select("created_at").First(&job, jobID).Error; err != nil {
		return time.Now()
	}
	return job.CreatedAt
}

// 輔助方法：轉換為響應 VO
func (s *QuantumJobService) convertToResponse(job *model.QuantumJob) vo.QuantumJobResponse {
	return vo.QuantumJobResponse{
		ID:                  job.ID,
		Title:               job.Title,
		Description:         job.Description,
		Status:              string(job.Status),
		Priority:            job.Priority,
		CreatedAt:           job.CreatedAt,
		UpdatedAt:           job.UpdatedAt,
		StartedAt:           job.StartedAt,
		CompletedAt:         job.CompletedAt,
		ExecutionTimeSeconds: job.ExecutionTimeSeconds,
		QuantumBackend:      job.QuantumBackend,
		IsSimulation:        job.IsSimulation,
		Tags:                job.Tags,
		ConfidenceScore:     job.ConfidenceScore,
		IsMalicious:         job.IsMalicious,
		Source:              job.Source,
		InputParamsSummary:  job.InputParams,
		ResultsSummary:      job.Results,
	}
}

// 輔助方法：轉換為詳細響應 VO
func (s *QuantumJobService) convertToDetailResponse(job *model.QuantumJob) *vo.QuantumJobDetailResponse {
	response := &vo.QuantumJobDetailResponse{
		QuantumJobResponse: s.convertToResponse(job),
		InputParams:        job.InputParams,
		Results:            job.Results,
		ErrorMessage:       job.ErrorMessage,
		Notes:              job.Notes,
		Logs:               make([]vo.QuantumJobLogResponse, len(job.Logs)),
	}

	// 轉換日誌
	for i, log := range job.Logs {
		response.Logs[i] = vo.QuantumJobLogResponse{
			ID:        log.ID,
			OldStatus: string(log.OldStatus),
			NewStatus: string(log.NewStatus),
			Message:   log.Message,
			CreatedBy: log.CreatedBy,
			CreatedAt: log.CreatedAt,
		}
	}

	return response
}
