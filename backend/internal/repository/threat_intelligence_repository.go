package repository

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/model"
)

// ThreatIntelligenceRepository 威脅情報儲存庫介面
type ThreatIntelligenceRepository interface {
	Create(ctx context.Context, threat *model.ThreatIntelligence) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.ThreatIntelligence, error)
	GetByIP(ctx context.Context, ip net.IP) ([]*model.ThreatIntelligence, error)
	GetByDomain(ctx context.Context, domain string) ([]*model.ThreatIntelligence, error)
	Update(ctx context.Context, threat *model.ThreatIntelligence) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *ThreatIntelligenceFilter) ([]*model.ThreatIntelligence, int64, error)
	BulkCreate(ctx context.Context, threats []*model.ThreatIntelligence) error
	GetStats(ctx context.Context, filter *StatsFilter) (*ThreatIntelligenceStats, error)
	GetRecentThreats(ctx context.Context, hours int, limit int) ([]*model.ThreatIntelligence, error)
	GetHighRiskThreats(ctx context.Context, limit int) ([]*model.ThreatIntelligence, error)
}

// ThreatIntelligenceFilter 威脅情報篩選器
type ThreatIntelligenceFilter struct {
	IPAddress   *string
	Domain      *string
	ThreatType  *string
	Severity    *string
	Source      *string
	CountryCode *string
	Tags        []string
	StartTime   *time.Time
	EndTime     *time.Time
	Page        int
	PageSize    int
	SortBy      string
	SortOrder   string
}

// StatsFilter 統計篩選器
type StatsFilter struct {
	GroupBy   string
	StartTime *time.Time
	EndTime   *time.Time
}

// ThreatIntelligenceStats 威脅情報統計
type ThreatIntelligenceStats struct {
	TotalThreats    int64
	HighRiskCount   int64
	RecentCount     int64
	CountByType     map[string]int64
	CountBySeverity map[string]int64
	CountBySource   map[string]int64
	CountByCountry  map[string]int64
	Timeline        []TimelinePoint
}

// TimelinePoint 時間線點
type TimelinePoint struct {
	Date  string
	Count int64
}

// threatIntelligenceRepository 威脅情報儲存庫實作
type threatIntelligenceRepository struct {
	db *gorm.DB
}

// NewThreatIntelligenceRepository 建立威脅情報儲存庫
func NewThreatIntelligenceRepository(db *gorm.DB) ThreatIntelligenceRepository {
	return &threatIntelligenceRepository{db: db}
}

// Create 建立威脅情報
func (r *threatIntelligenceRepository) Create(ctx context.Context, threat *model.ThreatIntelligence) error {
	return r.db.WithContext(ctx).Create(threat).Error
}

// GetByID 根據 ID 取得威脅情報
func (r *threatIntelligenceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ThreatIntelligence, error) {
	var threat model.ThreatIntelligence
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&threat).Error
	if err != nil {
		return nil, err
	}
	return &threat, nil
}

// GetByIP 根據 IP 取得威脅情報
func (r *threatIntelligenceRepository) GetByIP(ctx context.Context, ip net.IP) ([]*model.ThreatIntelligence, error) {
	var threats []*model.ThreatIntelligence
	err := r.db.WithContext(ctx).Where("ip_address = ?", ip).Find(&threats).Error
	return threats, err
}

// GetByDomain 根據域名取得威脅情報
func (r *threatIntelligenceRepository) GetByDomain(ctx context.Context, domain string) ([]*model.ThreatIntelligence, error) {
	var threats []*model.ThreatIntelligence
	err := r.db.WithContext(ctx).Where("domain = ?", domain).Find(&threats).Error
	return threats, err
}

// Update 更新威脅情報
func (r *threatIntelligenceRepository) Update(ctx context.Context, threat *model.ThreatIntelligence) error {
	return r.db.WithContext(ctx).Save(threat).Error
}

// Delete 刪除威脅情報
func (r *threatIntelligenceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.ThreatIntelligence{}, id).Error
}

// List 取得威脅情報列表
func (r *threatIntelligenceRepository) List(ctx context.Context, filter *ThreatIntelligenceFilter) ([]*model.ThreatIntelligence, int64, error) {
	var threats []*model.ThreatIntelligence
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ThreatIntelligence{})

	// 應用篩選器
	query = r.applyFilter(query, filter)

	// 計算總數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 應用分頁和排序
	if filter.SortBy != "" {
		order := fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder)
		query = query.Order(order)
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err = query.Find(&threats).Error
	return threats, total, err
}

// BulkCreate 批量建立威脅情報
func (r *threatIntelligenceRepository) BulkCreate(ctx context.Context, threats []*model.ThreatIntelligence) error {
	return r.db.WithContext(ctx).CreateInBatches(threats, 100).Error
}

// GetStats 取得統計資料
func (r *threatIntelligenceRepository) GetStats(ctx context.Context, filter *StatsFilter) (*ThreatIntelligenceStats, error) {
	stats := &ThreatIntelligenceStats{
		CountByType:     make(map[string]int64),
		CountBySeverity: make(map[string]int64),
		CountBySource:   make(map[string]int64),
		CountByCountry:  make(map[string]int64),
	}

	query := r.db.WithContext(ctx).Model(&model.ThreatIntelligence{})
	
	// 應用時間篩選
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", filter.EndTime)
	}

	// 總威脅數
	query.Count(&stats.TotalThreats)

	// 高風險威脅數
	query.Where("severity IN ?", []string{"high", "critical"}).Count(&stats.HighRiskCount)

	// 最近威脅數（24小時內）
	query.Where("last_seen >= ?", time.Now().Add(-24*time.Hour)).Count(&stats.RecentCount)

	// 按類型統計
	r.getCountByField(ctx, query, "threat_type", stats.CountByType)
	
	// 按嚴重程度統計
	r.getCountByField(ctx, query, "severity", stats.CountBySeverity)
	
	// 按來源統計
	r.getCountByField(ctx, query, "source", stats.CountBySource)
	
	// 按國家統計
	r.getCountByField(ctx, query, "country_code", stats.CountByCountry)

	return stats, nil
}

// GetRecentThreats 取得最近威脅
func (r *threatIntelligenceRepository) GetRecentThreats(ctx context.Context, hours int, limit int) ([]*model.ThreatIntelligence, error) {
	var threats []*model.ThreatIntelligence
	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	
	err := r.db.WithContext(ctx).
		Where("last_seen >= ?", since).
		Order("last_seen DESC").
		Limit(limit).
		Find(&threats).Error
	
	return threats, err
}

// GetHighRiskThreats 取得高風險威脅
func (r *threatIntelligenceRepository) GetHighRiskThreats(ctx context.Context, limit int) ([]*model.ThreatIntelligence, error) {
	var threats []*model.ThreatIntelligence
	
	err := r.db.WithContext(ctx).
		Where("severity IN ?", []string{"high", "critical"}).
		Order("confidence_score DESC, last_seen DESC").
		Limit(limit).
		Find(&threats).Error
	
	return threats, err
}

// applyFilter 應用篩選器
func (r *threatIntelligenceRepository) applyFilter(query *gorm.DB, filter *ThreatIntelligenceFilter) *gorm.DB {
	if filter.IPAddress != nil {
		query = query.Where("ip_address = ?", *filter.IPAddress)
	}
	if filter.Domain != nil {
		query = query.Where("domain = ?", *filter.Domain)
	}
	if filter.ThreatType != nil {
		query = query.Where("threat_type = ?", *filter.ThreatType)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}
	if filter.Source != nil {
		query = query.Where("source = ?", *filter.Source)
	}
	if filter.CountryCode != nil {
		query = query.Where("country_code = ?", *filter.CountryCode)
	}
	if len(filter.Tags) > 0 {
		query = query.Where("tags && ?", filter.Tags)
	}
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", filter.EndTime)
	}
	return query
}

// getCountByField 根據欄位取得統計數據
func (r *threatIntelligenceRepository) getCountByField(ctx context.Context, query *gorm.DB, field string, result map[string]int64) {
	var counts []struct {
		Field string
		Count int64
	}
	
	query.Select(fmt.Sprintf("%s as field, count(*) as count", field)).
		Group(field).
		Scan(&counts)
	
	for _, count := range counts {
		result[count.Field] = count.Count
	}
} 