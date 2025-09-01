package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/api/proto/api/proto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/dto"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// ThreatIntelligenceGRPCServer gRPC服務器實作
type ThreatIntelligenceGRPCServer struct {
	proto.UnimplementedThreatIntelligenceServiceServer
	threatService service.ThreatIntelligenceService
}

// NewThreatIntelligenceGRPCServer 建立gRPC服務器
func NewThreatIntelligenceGRPCServer(threatService service.ThreatIntelligenceService) *ThreatIntelligenceGRPCServer {
	return &ThreatIntelligenceGRPCServer{
		threatService: threatService,
	}
}

// GetThreatIntelligence 取得威脅情報
func (s *ThreatIntelligenceGRPCServer) GetThreatIntelligence(ctx context.Context, req *proto.GetThreatIntelligenceRequest) (*proto.GetThreatIntelligenceResponse, error) {
	// 解析UUID
	threatID, err := uuid.Parse(req.Id)
	if err != nil {
		return &proto.GetThreatIntelligenceResponse{
			Success:   false,
			Message:   "Invalid threat ID format",
			ErrorCode: "INVALID_ID",
		}, status.Error(codes.InvalidArgument, "Invalid threat ID format")
	}

	// 呼叫服務
	threat, err := s.threatService.GetByID(threatID)
	if err != nil {
		pkglogger.Error("Failed to get threat intelligence", pkglogger.Fields{
			"error":      err.Error(),
			"threat_id":  req.Id,
		})

		if err == dto.ErrThreatNotFound {
			return &proto.GetThreatIntelligenceResponse{
				Success:   false,
				Message:   "Threat intelligence not found",
				ErrorCode: "NOT_FOUND",
			}, status.Error(codes.NotFound, "Threat intelligence not found")
		}

		return &proto.GetThreatIntelligenceResponse{
			Success:   false,
			Message:   "Failed to get threat intelligence",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	// 轉換為proto格式
	protoThreat := convertThreatToProto(threat)

	return &proto.GetThreatIntelligenceResponse{
		Success: true,
		Message: "Threat intelligence retrieved successfully",
		Data:    protoThreat,
	}, nil
}

// ListThreatIntelligence 列出威脅情報
func (s *ThreatIntelligenceGRPCServer) ListThreatIntelligence(ctx context.Context, req *proto.ListThreatIntelligenceRequest) (*proto.ListThreatIntelligenceResponse, error) {
	// 轉換請求參數
	filter := dto.ThreatIntelligenceFilter{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		ThreatType: req.ThreatType,
		Severity:   req.Severity,
		Source:     req.Source,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		IsActive:   &req.IsActive,
	}

	// 設定預設值
	filter.SetDefaults()

	// 呼叫服務
	result, err := s.threatService.List(&filter)
	if err != nil {
		pkglogger.Error("Failed to list threat intelligence", pkglogger.Fields{
			"error": err.Error(),
		})

		return &proto.ListThreatIntelligenceResponse{
			Success:   false,
			Message:   "Failed to list threat intelligence",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	// 轉換為proto格式
	var protoThreats []*proto.ThreatIntelligence
	for _, threat := range result.Threats {
		protoThreats = append(protoThreats, convertThreatToProto(&threat))
	}

	pagination := &proto.Pagination{
		CurrentPage:  int32(result.Pagination.CurrentPage),
		PageSize:     int32(result.Pagination.PageSize),
		TotalPages:   int32(result.Pagination.TotalPages),
		TotalRecords: result.Pagination.TotalRecords,
		HasNext:      result.Pagination.HasNext,
		HasPrevious:  result.Pagination.HasPrevious,
	}

	return &proto.ListThreatIntelligenceResponse{
		Success:    true,
		Message:    "Threat intelligence list retrieved successfully",
		Threats:    protoThreats,
		Pagination: pagination,
	}, nil
}

// CreateThreatIntelligence 建立威脅情報
func (s *ThreatIntelligenceGRPCServer) CreateThreatIntelligence(ctx context.Context, req *proto.CreateThreatIntelligenceRequest) (*proto.CreateThreatIntelligenceResponse, error) {
	// 轉換請求為DTO
	createReq := &dto.CreateThreatIntelligenceRequest{
		IPAddress:       req.IpAddress,
		Domain:          req.Domain,
		URL:             req.Url,
		FileHash:        req.FileHash,
		ThreatType:      req.ThreatType,
		Severity:        req.Severity,
		Description:     req.Description,
		Tags:            req.Tags,
		ConfidenceScore: int(req.ConfidenceScore),
		Source:          req.Source,
	}

	// 轉換時間戳
	if req.FirstSeen != nil {
		firstSeen := req.FirstSeen.AsTime()
		createReq.FirstSeen = &firstSeen
	}
	if req.LastSeen != nil {
		lastSeen := req.LastSeen.AsTime()
		createReq.LastSeen = &lastSeen
	}

	// 轉換metadata
	if req.Metadata != nil {
		createReq.Metadata = req.Metadata
	}

	// 呼叫服務
	threat, err := s.threatService.Create(createReq)
	if err != nil {
		pkglogger.Error("Failed to create threat intelligence", pkglogger.Fields{
			"error": err.Error(),
		})

		return &proto.CreateThreatIntelligenceResponse{
			Success:   false,
			Message:   "Failed to create threat intelligence",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	// 轉換為proto格式
	protoThreat := convertThreatToProto(threat)

	return &proto.CreateThreatIntelligenceResponse{
		Success: true,
		Message: "Threat intelligence created successfully",
		Data:    protoThreat,
	}, nil
}

// UpdateThreatIntelligence 更新威脅情報
func (s *ThreatIntelligenceGRPCServer) UpdateThreatIntelligence(ctx context.Context, req *proto.UpdateThreatIntelligenceRequest) (*proto.UpdateThreatIntelligenceResponse, error) {
	// 解析UUID
	threatID, err := uuid.Parse(req.Id)
	if err != nil {
		return &proto.UpdateThreatIntelligenceResponse{
			Success:   false,
			Message:   "Invalid threat ID format",
			ErrorCode: "INVALID_ID",
		}, status.Error(codes.InvalidArgument, "Invalid threat ID format")
	}

	// 轉換請求為DTO
	updateReq := &dto.UpdateThreatIntelligenceRequest{
		IPAddress:       &req.IpAddress,
		Domain:          &req.Domain,
		URL:             &req.Url,
		FileHash:        &req.FileHash,
		ThreatType:      &req.ThreatType,
		Severity:        &req.Severity,
		Description:     &req.Description,
		Tags:            req.Tags,
		ConfidenceScore: func() *int { cs := int(req.ConfidenceScore); return &cs }(),
		Source:          &req.Source,
		IsActive:        &req.IsActive,
	}

	// 轉換時間戳
	if req.FirstSeen != nil {
		firstSeen := req.FirstSeen.AsTime()
		updateReq.FirstSeen = &firstSeen
	}
	if req.LastSeen != nil {
		lastSeen := req.LastSeen.AsTime()
		updateReq.LastSeen = &lastSeen
	}

	// 轉換metadata
	if req.Metadata != nil {
		updateReq.Metadata = req.Metadata
	}

	// 呼叫服務
	threat, err := s.threatService.Update(threatID, updateReq)
	if err != nil {
		pkglogger.Error("Failed to update threat intelligence", pkglogger.Fields{
			"error":     err.Error(),
			"threat_id": req.Id,
		})

		if err == dto.ErrThreatNotFound {
			return &proto.UpdateThreatIntelligenceResponse{
				Success:   false,
				Message:   "Threat intelligence not found",
				ErrorCode: "NOT_FOUND",
			}, status.Error(codes.NotFound, "Threat intelligence not found")
		}

		return &proto.UpdateThreatIntelligenceResponse{
			Success:   false,
			Message:   "Failed to update threat intelligence",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	// 轉換為proto格式
	protoThreat := convertThreatToProto(threat)

	return &proto.UpdateThreatIntelligenceResponse{
		Success: true,
		Message: "Threat intelligence updated successfully",
		Data:    protoThreat,
	}, nil
}

// DeleteThreatIntelligence 刪除威脅情報
func (s *ThreatIntelligenceGRPCServer) DeleteThreatIntelligence(ctx context.Context, req *proto.DeleteThreatIntelligenceRequest) (*proto.DeleteThreatIntelligenceResponse, error) {
	// 解析UUID
	threatID, err := uuid.Parse(req.Id)
	if err != nil {
		return &proto.DeleteThreatIntelligenceResponse{
			Success:   false,
			Message:   "Invalid threat ID format",
			ErrorCode: "INVALID_ID",
		}, status.Error(codes.InvalidArgument, "Invalid threat ID format")
	}

	// 呼叫服務
	err = s.threatService.Delete(threatID)
	if err != nil {
		pkglogger.Error("Failed to delete threat intelligence", pkglogger.Fields{
			"error":     err.Error(),
			"threat_id": req.Id,
		})

		if err == dto.ErrThreatNotFound {
			return &proto.DeleteThreatIntelligenceResponse{
				Success:   false,
				Message:   "Threat intelligence not found",
				ErrorCode: "NOT_FOUND",
			}, status.Error(codes.NotFound, "Threat intelligence not found")
		}

		return &proto.DeleteThreatIntelligenceResponse{
			Success:   false,
			Message:   "Failed to delete threat intelligence",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	return &proto.DeleteThreatIntelligenceResponse{
		Success: true,
		Message: "Threat intelligence deleted successfully",
	}, nil
}

// SearchThreatIntelligence 搜尋威脅情報
func (s *ThreatIntelligenceGRPCServer) SearchThreatIntelligence(ctx context.Context, req *proto.SearchThreatIntelligenceRequest) (*proto.SearchThreatIntelligenceResponse, error) {
	// 轉換請求參數
	searchReq := &dto.SearchThreatIntelligenceRequest{
		Query:       req.Query,
		IPAddress:   req.IpAddress,
		Domain:      req.Domain,
		ThreatType:  req.ThreatType,
		Severity:    req.Severity,
		Page:        int(req.Page),
		PageSize:    int(req.PageSize),
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}

	// 設定預設值
	searchReq.SetDefaults()

	// 呼叫服務
	result, err := s.threatService.Search(searchReq)
	if err != nil {
		pkglogger.Error("Failed to search threat intelligence", pkglogger.Fields{
			"error": err.Error(),
			"query": req.Query,
		})

		return &proto.SearchThreatIntelligenceResponse{
			Success:   false,
			Message:   "Failed to search threat intelligence",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	// 轉換為proto格式
	var protoThreats []*proto.ThreatIntelligence
	for _, threat := range result.Threats {
		protoThreats = append(protoThreats, convertThreatToProto(&threat))
	}

	pagination := &proto.Pagination{
		CurrentPage:  int32(result.Pagination.CurrentPage),
		PageSize:     int32(result.Pagination.PageSize),
		TotalPages:   int32(result.Pagination.TotalPages),
		TotalRecords: result.Pagination.TotalRecords,
		HasNext:      result.Pagination.HasNext,
		HasPrevious:  result.Pagination.HasPrevious,
	}

	return &proto.SearchThreatIntelligenceResponse{
		Success:    true,
		Message:    "Search completed successfully",
		Threats:    protoThreats,
		Pagination: pagination,
	}, nil
}

// GetThreatStatistics 取得統計資訊
func (s *ThreatIntelligenceGRPCServer) GetThreatStatistics(ctx context.Context, req *proto.GetThreatStatisticsRequest) (*proto.GetThreatStatisticsResponse, error) {
	// 呼叫服務取得統計資訊
	stats, err := s.threatService.GetStatistics()
	if err != nil {
		pkglogger.Error("Failed to get threat statistics", pkglogger.Fields{
			"error": err.Error(),
		})

		return &proto.GetThreatStatisticsResponse{
			Success:   false,
			Message:   "Failed to get threat statistics",
			ErrorCode: "INTERNAL_ERROR",
		}, status.Error(codes.Internal, "Internal server error")
	}

	// 轉換為proto格式
	protoStats := &proto.ThreatStatistics{
		TotalThreats:     stats.TotalThreats,
		ActiveThreats:    stats.ActiveThreats,
		HighRiskThreats:  stats.HighRiskThreats,
		MediumRiskThreats: stats.MediumRiskThreats,
		LowRiskThreats:   stats.LowRiskThreats,
		ThreatTypes:      make(map[string]int64),
		Sources:          make(map[string]int64),
	}

	// 轉換威脅類型統計
	for _, stat := range stats.ThreatTypeStats {
		protoStats.ThreatTypes[stat.Label] = int64(stat.Value)
	}

	// 轉換來源統計
	for _, stat := range stats.SourceStats {
		protoStats.Sources[stat.Label] = int64(stat.Value)
	}

	return &proto.GetThreatStatisticsResponse{
		Success: true,
		Message: "Statistics retrieved successfully",
		Data:    protoStats,
	}, nil
}

// SubscribeThreats 即時威脅訂閱
func (s *ThreatIntelligenceGRPCServer) SubscribeThreats(req *proto.SubscribeThreatsRequest, stream proto.ThreatIntelligenceService_SubscribeThreatsServer) error {
	// TODO: 實作即時威脅訂閱功能
	// 這需要與MQTT或其他訊息佇列系統整合

	pkglogger.Info("Threat subscription started", pkglogger.Fields{
		"threat_types":         req.ThreatTypes,
		"severities":          req.Severities,
		"min_confidence_score": req.MinConfidenceScore,
	})

	// 模擬訂閱，實際實作中應該連接到訊息佇列
	for {
		select {
		case <-stream.Context().Done():
			pkglogger.Info("Threat subscription ended")
			return nil
		case <-time.After(30 * time.Second):
			// 發送心跳通知
			notification := &proto.ThreatNotification{
				NotificationId: uuid.New().String(),
				Type:          "heartbeat",
				Timestamp:     timestamppb.Now(),
				Metadata: map[string]string{
					"message": "Connection alive",
				},
			}

			if err := stream.Send(notification); err != nil {
				pkglogger.Error("Failed to send threat notification", pkglogger.Fields{
					"error": err.Error(),
				})
				return err
			}
		}
	}
}

// convertThreatToProto 轉換威脅情報到proto格式
func convertThreatToProto(threat *vo.ThreatIntelligenceVO) *proto.ThreatIntelligence {
	protoThreat := &proto.ThreatIntelligence{
		Id:              threat.ID.String(),
		IpAddress:       threat.IPAddress,
		Domain:          threat.Domain,
		Url:             threat.URL,
		FileHash:        threat.FileHash,
		ThreatType:      threat.ThreatType,
		Severity:        threat.Severity,
		Description:     threat.Description,
		Tags:            threat.Tags,
		ConfidenceScore: int32(threat.ConfidenceScore),
		Source:          threat.Source,
		IsActive:        threat.IsActive,
		RiskScore:       threat.RiskScore,
	}

	// 轉換時間戳
	if !threat.FirstSeen.IsZero() {
		protoThreat.FirstSeen = timestamppb.New(threat.FirstSeen)
	}
	if !threat.LastSeen.IsZero() {
		protoThreat.LastSeen = timestamppb.New(threat.LastSeen)
	}
	if !threat.CreatedAt.IsZero() {
		protoThreat.CreatedAt = timestamppb.New(threat.CreatedAt)
	}
	if !threat.UpdatedAt.IsZero() {
		protoThreat.UpdatedAt = timestamppb.New(threat.UpdatedAt)
	}

	// 轉換metadata
	if threat.Metadata != nil {
		protoThreat.Metadata = make(map[string]string)
		for k, v := range threat.Metadata {
			if str, ok := v.(string); ok {
				protoThreat.Metadata[k] = str
			} else {
				protoThreat.Metadata[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	return protoThreat
} 