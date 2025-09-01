package service

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/model"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/repository"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
)

// ThreatIntelligenceService 威脅情報服務介面
type ThreatIntelligenceService interface {
	CreateThreat(ctx context.Context, req *dto.ThreatIntelligenceCreateRequest) (*vo.ThreatIntelligenceVO, error)
	GetThreatByID(ctx context.Context, id uuid.UUID) (*vo.ThreatIntelligenceVO, error)
	UpdateThreat(ctx context.Context, id uuid.UUID, req *dto.ThreatIntelligenceUpdateRequest) (*vo.ThreatIntelligenceVO, error)
	DeleteThreat(ctx context.Context, id uuid.UUID) error
	ListThreats(ctx context.Context, req *dto.ThreatIntelligenceQueryRequest) (*vo.ThreatIntelligenceListVO, error)
	LookupIP(ctx context.Context, req *dto.ThreatIntelligenceIPLookupRequest) (*vo.ThreatIntelligenceIPLookupVO, error)
	LookupDomain(ctx context.Context, req *dto.ThreatIntelligenceDomainLookupRequest) (*vo.ThreatIntelligenceDomainLookupVO, error)
	GetStats(ctx context.Context, req *dto.ThreatIntelligenceStatsRequest) (*vo.ThreatIntelligenceStatsVO, error)
	BulkCreateThreats(ctx context.Context, req *dto.ThreatIntelligenceBulkCreateRequest) (*vo.ThreatIntelligenceBulkCreateVO, error)
	SearchThreats(ctx context.Context, query string, page, limit int) ([]*vo.ThreatIntelligenceVO, int64, error)
	GetStatistics(ctx context.Context) (map[string]interface{}, error)
	BulkUpdateThreats(ctx context.Context, req *dto.ThreatIntelligenceBulkUpdateRequest) (*vo.ThreatIntelligenceBulkUpdateVO, error)
	BulkDeleteThreats(ctx context.Context, req *dto.ThreatIntelligenceBulkDeleteRequest) (*vo.ThreatIntelligenceBulkDeleteVO, error)
}

// threatIntelligenceService 威脅情報服務實作
type threatIntelligenceService struct {
	repo repository.ThreatIntelligenceRepository
}

// NewThreatIntelligenceService 建立威脅情報服務
func NewThreatIntelligenceService(repo repository.ThreatIntelligenceRepository) ThreatIntelligenceService {
	return &threatIntelligenceService{repo: repo}
}

// CreateThreat 建立威脅情報
func (s *threatIntelligenceService) CreateThreat(ctx context.Context, req *dto.ThreatIntelligenceCreateRequest) (*vo.ThreatIntelligenceVO, error) {
	// 驗證 IP 地址
	ip := net.ParseIP(req.IPAddress)
	if ip == nil {
		return nil, dto.ErrInvalidIPAddress
	}

	// 建立威脅情報模型
	threat := &model.ThreatIntelligence{
		IPAddress:       ip,
		ThreatType:      model.ThreatType(req.ThreatType),
		Severity:        model.SeverityLevel(req.Severity),
		ConfidenceScore: req.ConfidenceScore,
		Source:          req.Source,
	}

	// 複製其他欄位
	if err := copier.Copy(threat, req); err != nil {
		return nil, err
	}

	// 設定標籤
	threat.Tags = model.StringArray(req.Tags)

	// 設定元資料
	if req.Metadata != nil {
		threat.Metadata = model.JSONB(req.Metadata)
	}

	// 建立威脅情報
	if err := s.repo.Create(ctx, threat); err != nil {
		return nil, err
	}

	// 轉換為 VO
	return s.modelToVO(threat), nil
}

// GetThreatByID 根據 ID 取得威脅情報
func (s *threatIntelligenceService) GetThreatByID(ctx context.Context, id uuid.UUID) (*vo.ThreatIntelligenceVO, error) {
	threat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.modelToVO(threat), nil
}

// UpdateThreat 更新威脅情報
func (s *threatIntelligenceService) UpdateThreat(ctx context.Context, id uuid.UUID, req *dto.ThreatIntelligenceUpdateRequest) (*vo.ThreatIntelligenceVO, error) {
	// 取得現有威脅情報
	threat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新欄位
	if req.ThreatType != nil {
		threat.ThreatType = model.ThreatType(*req.ThreatType)
	}
	if req.Severity != nil {
		threat.Severity = model.SeverityLevel(*req.Severity)
	}
	if req.ConfidenceScore != nil {
		threat.ConfidenceScore = *req.ConfidenceScore
	}
	if req.Description != nil {
		threat.Description = req.Description
	}
	if req.CountryCode != nil {
		threat.CountryCode = req.CountryCode
	}
	if req.ASN != nil {
		threat.ASN = req.ASN
	}
	if req.ISP != nil {
		threat.ISP = req.ISP
	}
	if req.Tags != nil {
		threat.Tags = model.StringArray(req.Tags)
	}
	if req.Metadata != nil {
		threat.Metadata = model.JSONB(req.Metadata)
	}

	// 更新威脅情報
	if err := s.repo.Update(ctx, threat); err != nil {
		return nil, err
	}

	return s.modelToVO(threat), nil
}

// DeleteThreat 刪除威脅情報
func (s *threatIntelligenceService) DeleteThreat(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ListThreats 取得威脅情報列表
func (s *threatIntelligenceService) ListThreats(ctx context.Context, req *dto.ThreatIntelligenceQueryRequest) (*vo.ThreatIntelligenceListVO, error) {
	// 設定預設值
	req.SetDefaults()

	// 建立篩選器
	filter := &repository.ThreatIntelligenceFilter{
		IPAddress:   req.IPAddress,
		Domain:      req.Domain,
		ThreatType:  req.ThreatType,
		Severity:    req.Severity,
		Source:      req.Source,
		CountryCode: req.CountryCode,
		Tags:        req.Tags,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Page:        req.Page,
		PageSize:    req.PageSize,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}

	// 取得威脅情報列表
	threats, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 轉換為 VO
	threatVOs := make([]vo.ThreatIntelligenceVO, len(threats))
	for i, threat := range threats {
		threatVOs[i] = *s.modelToVO(threat)
	}

	// 建立分頁資訊
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	pagination := vo.PaginationVO{
		CurrentPage:  req.Page,
		PageSize:     req.PageSize,
		TotalPages:   totalPages,
		TotalRecords: total,
		HasNext:      req.Page < totalPages,
		HasPrevious:  req.Page > 1,
	}

	return &vo.ThreatIntelligenceListVO{
		Data:       threatVOs,
		Pagination: pagination,
	}, nil
}

// LookupIP IP 查詢
func (s *threatIntelligenceService) LookupIP(ctx context.Context, req *dto.ThreatIntelligenceIPLookupRequest) (*vo.ThreatIntelligenceIPLookupVO, error) {
	ip := net.ParseIP(req.IPAddress)
	if ip == nil {
		return nil, dto.ErrInvalidIPAddress
	}

	threats, err := s.repo.GetByIP(ctx, ip)
	if err != nil {
		return nil, err
	}

	result := &vo.ThreatIntelligenceIPLookupVO{
		IPAddress:     req.IPAddress,
		IsKnownThreat: len(threats) > 0,
		ThreatCount:   len(threats),
		Details:       make([]vo.ThreatIntelligenceVO, len(threats)),
	}

	if len(threats) > 0 {
		// 找最高嚴重程度
		highestSeverity := ""
		severityOrder := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
		maxOrder := 0

		sources := make(map[string]bool)
		
		for i, threat := range threats {
			result.Details[i] = *s.modelToVO(threat)
			
			// 記錄來源
			sources[threat.Source] = true
			
			// 找最高嚴重程度
			if order := severityOrder[string(threat.Severity)]; order > maxOrder {
				maxOrder = order
				highestSeverity = string(threat.Severity)
			}
			
			// 設定時間
			if result.FirstSeen == nil || threat.FirstSeen.Before(*result.FirstSeen) {
				result.FirstSeen = &threat.FirstSeen
			}
			if result.LastSeen == nil || threat.LastSeen.After(*result.LastSeen) {
				result.LastSeen = &threat.LastSeen
			}
		}

		result.HighestSeverity = highestSeverity
		
		// 建立來源列表
		for source := range sources {
			result.Sources = append(result.Sources, source)
		}
	}

	return result, nil
}

// LookupDomain 域名查詢
func (s *threatIntelligenceService) LookupDomain(ctx context.Context, req *dto.ThreatIntelligenceDomainLookupRequest) (*vo.ThreatIntelligenceDomainLookupVO, error) {
	threats, err := s.repo.GetByDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}

	result := &vo.ThreatIntelligenceDomainLookupVO{
		Domain:        req.Domain,
		IsKnownThreat: len(threats) > 0,
		ThreatCount:   len(threats),
		Details:       make([]vo.ThreatIntelligenceVO, len(threats)),
	}

	if len(threats) > 0 {
		// 類似 IP 查詢的邏輯
		highestSeverity := ""
		severityOrder := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
		maxOrder := 0

		sources := make(map[string]bool)
		
		for i, threat := range threats {
			result.Details[i] = *s.modelToVO(threat)
			sources[threat.Source] = true
			
			if order := severityOrder[string(threat.Severity)]; order > maxOrder {
				maxOrder = order
				highestSeverity = string(threat.Severity)
			}
			
			if result.FirstSeen == nil || threat.FirstSeen.Before(*result.FirstSeen) {
				result.FirstSeen = &threat.FirstSeen
			}
			if result.LastSeen == nil || threat.LastSeen.After(*result.LastSeen) {
				result.LastSeen = &threat.LastSeen
			}
		}

		result.HighestSeverity = highestSeverity
		for source := range sources {
			result.Sources = append(result.Sources, source)
		}
	}

	return result, nil
}

// GetStats 取得統計資料
func (s *threatIntelligenceService) GetStats(ctx context.Context, req *dto.ThreatIntelligenceStatsRequest) (*vo.ThreatIntelligenceStatsVO, error) {
	filter := &repository.StatsFilter{
		GroupBy:   req.GroupBy,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	stats, err := s.repo.GetStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &vo.ThreatIntelligenceStatsVO{
		TotalThreats:    int(stats.TotalThreats),
		HighRiskCount:   int(stats.HighRiskCount),
		RecentCount:     int(stats.RecentCount),
		CountByType:     s.convertCountMap(stats.CountByType),
		CountBySeverity: s.convertCountMap(stats.CountBySeverity),
		CountBySource:   s.convertCountMap(stats.CountBySource),
		CountByCountry:  s.convertCountMap(stats.CountByCountry),
	}, nil
}

// BulkCreateThreats 批量建立威脅情報
func (s *threatIntelligenceService) BulkCreateThreats(ctx context.Context, req *dto.ThreatIntelligenceBulkCreateRequest) (*vo.ThreatIntelligenceBulkCreateVO, error) {
	var successThreats []vo.ThreatIntelligenceVO
	var failedErrors []vo.BulkOperationError

	threats := make([]*model.ThreatIntelligence, 0, len(req.Items))

	for i, item := range req.Items {
		// 驗證並轉換每個項目
		ip := net.ParseIP(item.IPAddress)
		if ip == nil {
			failedErrors = append(failedErrors, vo.BulkOperationError{
				Index:   i,
				Error:   "INVALID_IP",
				Message: "Invalid IP address: " + item.IPAddress,
			})
			continue
		}

		threat := &model.ThreatIntelligence{
			IPAddress:       ip,
			ThreatType:      model.ThreatType(item.ThreatType),
			Severity:        model.SeverityLevel(item.Severity),
			ConfidenceScore: item.ConfidenceScore,
			Source:          item.Source,
		}

		if err := copier.Copy(threat, &item); err != nil {
			failedErrors = append(failedErrors, vo.BulkOperationError{
				Index:   i,
				Error:   "COPY_ERROR",
				Message: err.Error(),
			})
			continue
		}

		threat.Tags = model.StringArray(item.Tags)
		if item.Metadata != nil {
			threat.Metadata = model.JSONB(item.Metadata)
		}

		threats = append(threats, threat)
	}

	// 批量建立威脅情報
	if len(threats) > 0 {
		if err := s.repo.BulkCreate(ctx, threats); err != nil {
			return nil, err
		}

		// 轉換成功的項目
		for _, threat := range threats {
			successThreats = append(successThreats, *s.modelToVO(threat))
		}
	}

	return &vo.ThreatIntelligenceBulkCreateVO{
		Success:      successThreats,
		Failed:       failedErrors,
		TotalCount:   len(req.Items),
		SuccessCount: len(successThreats),
		FailedCount:  len(failedErrors),
	}, nil
}

// modelToVO 將模型轉換為 VO
func (s *threatIntelligenceService) modelToVO(threat *model.ThreatIntelligence) *vo.ThreatIntelligenceVO {
	threatVO := &vo.ThreatIntelligenceVO{
		ID:              threat.ID,
		IPAddress:       threat.IPAddress.String(),
		ThreatType:      string(threat.ThreatType),
		Severity:        string(threat.Severity),
		ConfidenceScore: threat.ConfidenceScore,
		Source:          threat.Source,
		FirstSeen:       threat.FirstSeen,
		LastSeen:        threat.LastSeen,
		Tags:            []string(threat.Tags),
		Metadata:        map[string]interface{}(threat.Metadata),
		RiskScore:       threat.GetRiskScore(),
		IsHighRisk:      threat.IsHighRisk(),
		IsRecent:        threat.IsRecent(),
		CreatedAt:       threat.CreatedAt,
		UpdatedAt:       threat.UpdatedAt,
	}

	// 複製指標欄位
	if threat.Domain != nil {
		threatVO.Domain = threat.Domain
	}
	if threat.Description != nil {
		threatVO.Description = threat.Description
	}
	if threat.ExternalID != nil {
		threatVO.ExternalID = threat.ExternalID
	}
	if threat.CountryCode != nil {
		threatVO.CountryCode = threat.CountryCode
	}
	if threat.ASN != nil {
		threatVO.ASN = threat.ASN
	}
	if threat.ISP != nil {
		threatVO.ISP = threat.ISP
	}

	return threatVO
}

// convertCountMap 轉換計數對應表
func (s *threatIntelligenceService) convertCountMap(m map[string]int64) map[string]int {
	result := make(map[string]int)
	for k, v := range m {
		result[k] = int(v)
	}
	return result
}

// SearchThreats 搜尋威脅情報
func (s *threatIntelligenceService) SearchThreats(ctx context.Context, query string, page, limit int) ([]*vo.ThreatIntelligenceVO, int64, error) {
	// 實現搜尋邏輯
	// TODO: 實現實際的搜尋功能
	return []*vo.ThreatIntelligenceVO{}, 0, nil
}

// GetStatistics 取得統計資訊
func (s *threatIntelligenceService) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	// 實現統計邏輯
	// TODO: 實現實際的統計功能
	return map[string]interface{}{
		"total_threats": 0,
		"high_risk":     0,
		"recent":        0,
	}, nil
}

// BulkUpdateThreats 批量更新威脅情報
func (s *threatIntelligenceService) BulkUpdateThreats(ctx context.Context, req *dto.ThreatIntelligenceBulkUpdateRequest) (*vo.ThreatIntelligenceBulkUpdateVO, error) {
	// 實現批量更新邏輯
	// TODO: 實現實際的批量更新功能
	return &vo.ThreatIntelligenceBulkUpdateVO{
		Success:      []vo.ThreatIntelligenceVO{},
		Failed:       []vo.BulkOperationError{},
		TotalCount:   0,
		SuccessCount: 0,
		FailedCount:  0,
	}, nil
}

// BulkDeleteThreats 批量刪除威脅情報
func (s *threatIntelligenceService) BulkDeleteThreats(ctx context.Context, req *dto.ThreatIntelligenceBulkDeleteRequest) (*vo.ThreatIntelligenceBulkDeleteVO, error) {
	// 實現批量刪除邏輯
	// TODO: 實現實際的批量刪除功能
	return &vo.ThreatIntelligenceBulkDeleteVO{
		Success:      []string{},
		Failed:       []vo.BulkOperationError{},
		TotalCount:   0,
		SuccessCount: 0,
		FailedCount:  0,
	}, nil
} 