package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全域資料庫實例
var DB *gorm.DB

// Init 初始化資料庫連線
func Init(dsn string) (*gorm.DB, error) {
	// 配置 GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // 慢查詢閾值
			LogLevel:                  logger.Info,   // 日誌級別
			IgnoreRecordNotFoundError: true,          // 忽略 ErrRecordNotFound 錯誤
			Colorful:                  true,          // 彩色輸出
		},
	)

	// 連接到資料庫
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction: false,
	})
	
	if err != nil {
		return nil, fmt.Errorf("無法連接到資料庫: %w", err)
	}

	// 取得底層的 sql.DB 以配置連線池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("無法取得 SQL DB 實例: %w", err)
	}

	// 配置連線池
	sqlDB.SetMaxOpenConns(25)                 // 最大開啟連線數
	sqlDB.SetMaxIdleConns(25)                 // 最大閒置連線數
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 連線最大生命週期

	// 測試連線
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("資料庫連線測試失敗: %w", err)
	}

	// 設定全域變數
	DB = db

	log.Println("資料庫連線成功建立")
	return db, nil
}

// Close 關閉資料庫連線
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("無法取得 SQL DB 實例: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 取得資料庫實例
func GetDB() *gorm.DB {
	return DB
}

// HealthCheck 檢查資料庫健康狀態
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("資料庫未初始化")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("無法取得 SQL DB 實例: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("資料庫連線失敗: %w", err)
	}

	return nil
}

// GetStats 取得資料庫連線統計
func GetStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{
			"status": "not_initialized",
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"status":           "connected",
		"open_connections": stats.OpenConnections,
		"in_use":          stats.InUse,
		"idle":            stats.Idle,
		"wait_count":      stats.WaitCount,
		"wait_duration":   stats.WaitDuration,
		"max_idle_closed": stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
} 