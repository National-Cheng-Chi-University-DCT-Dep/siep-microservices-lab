package model

import (
	"gorm.io/gorm"
)

// AllModels 返回所有模型的切片，用於自動遷移
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&APIKey{},
		&ThreatIntelligence{},
		&IntelligenceSource{},
		&CollectionJob{},
	}
}

// AutoMigrate 自動遷移所有模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}

// Constants 匯出所有常數
const (
	// User Roles
	UserRoleAdmin   = RoleAdmin
	UserRolePremium = RolePremium
	UserRoleBasic   = RoleBasic

	// Threat Types
	ThreatTypeMalware    = ThreatMalware
	ThreatTypePhishing   = ThreatPhishing
	ThreatTypeSpam       = ThreatSpam
	ThreatTypeBotnet     = ThreatBotnet
	ThreatTypeScanner    = ThreatScanner
	ThreatTypeDDoS       = ThreatDDoS
	ThreatTypeBruteforce = ThreatBruteforce
	ThreatTypeOther      = ThreatOther

	// Severity Levels
	SeverityLevelLow      = SeverityLow
	SeverityLevelMedium   = SeverityMedium
	SeverityLevelHigh     = SeverityHigh
	SeverityLevelCritical = SeverityCritical

	// Collection Status
	CollectionStatusPending    = StatusPending
	CollectionStatusInProgress = StatusInProgress
	CollectionStatusCompleted  = StatusCompleted
	CollectionStatusFailed     = StatusFailed
) 