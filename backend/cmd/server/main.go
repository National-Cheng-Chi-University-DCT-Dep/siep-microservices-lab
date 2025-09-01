package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/collector"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/config"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/handler"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/middleware"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/model"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/repository"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/service"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/database"
	pkgjwt "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/jwt"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
	pkgmqtt "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/mqtt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 資安情報平台 API
// @version 1.0
// @description 一個創新、自動化且可擴展的開源資安威脅情報平台 API
// @termsOfService http://swagger.io/terms/

// @contact.name API 支援
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://github.com/lipeichen/Ultimate-Security-Intelligence-Platform/blob/main/LICENSE

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 檔案，使用系統環境變數")
	}

	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("無法載入配置:", err)
	}

	// 初始化日誌
	logger.Init(cfg.LogLevel)
	logger.Info("啟動資安情報平台", logger.Fields{
		"environment": cfg.Environment,
		"port":        cfg.Server.Port,
		"version":     "1.0.0",
	})

	// 初始化資料庫
	db, err := database.Init(cfg.Database.DSN)
	if err != nil {
		logger.Error("無法連接資料庫", logger.Fields{
			"error": err.Error(),
			"dsn":   cfg.Database.DSN,
		})
		log.Fatal("資料庫連接失敗")
	}
	defer database.Close()

	// 自動遷移資料庫模式
	if err := model.AutoMigrate(db); err != nil {
		logger.Error("資料庫遷移失敗", logger.Fields{
			"error": err.Error(),
		})
		log.Fatal("資料庫遷移失敗")
	}

	// 初始化JWT管理器
	jwtManager := pkgjwt.NewJWTManager(
		getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		getEnvOrDefault("JWT_ISSUER", "security-intelligence-platform"),
		24, // 24小時過期
	)

	// 初始化MQTT客戶端
	var mqttClient pkgmqtt.MQTTClientInterface
	
	if mqttBroker := getEnvOrDefault("MQTT_BROKER", ""); mqttBroker != "" {
		mqttClient = pkgmqtt.NewMQTTClient(
			mqttBroker,
			getEnvOrDefault("MQTT_USERNAME", ""),
			getEnvOrDefault("MQTT_PASSWORD", ""),
		)
		
		if err := mqttClient.Connect(); err != nil {
			logger.Warn("MQTT連接失敗，將在沒有實時通知的情況下運行", logger.Fields{
				"error": err.Error(),
			})
		} else {
			logger.Info("MQTT客戶端已連接", logger.Fields{
				"broker": mqttBroker,
			})
		}
	}

	// 初始化Repository層
	threatIntelRepo := repository.NewThreatIntelligenceRepository(db)

	// 初始化Service層
	threatIntelService := service.NewThreatIntelligenceService(threatIntelRepo)
	authService := service.NewAuthService(db, jwtManager)

	// 初始化威脅情報收集器
	_ = collector.NewAbuseIPDBCollector(
		getEnvOrDefault("ABUSEIPDB_API_KEY", ""),
		threatIntelService,
	)
	
	// 初始化 HIBP 收集器
	hibpCollector := collector.NewHIBPCollector(
		getEnvOrDefault("HIBP_API_KEY", ""),
		threatIntelService,
	)

	// 初始化Handler層
	threatIntelHandler := handler.NewThreatIntelligenceHandler(threatIntelService)
	collectorHandler := handler.NewCollectorHandler(threatIntelService)
	authHandler := handler.NewAuthHandler(authService)
	hibpHandler := handler.NewHIBPHandler(hibpCollector, threatIntelService)

	// 創建gRPC服務器
	// TODO: 修復 gRPC 服務器
	// grpcServer := grpc.NewServer()
	// threatGRPCHandler := grpchandler.NewThreatIntelligenceGRPCServer(threatIntelService)
	// proto.RegisterThreatIntelligenceServiceServer(grpcServer, threatGRPCHandler)

	// 設定Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 建立Gin路由器
	r := gin.New()

	// 設定中介軟體
	setupMiddlewares(r, cfg, jwtManager)

	// 設定路由
	setupRoutes(r, threatIntelHandler, collectorHandler, authHandler, hibpHandler, jwtManager)

	// 創建HTTP伺服器
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// 啟動gRPC伺服器
	// TODO: 修復 gRPC 服務器
	// grpcPort := getEnvOrDefault("GRPC_PORT", "9090")
	// go func() {
	// 	lis, err := net.Listen("tcp", ":"+grpcPort)
	// 	if err != nil {
	// 		logger.Error("gRPC伺服器監聽失敗", logger.Fields{
	// 			"error": err.Error(),
	// 			"port":  grpcPort,
	// 		})
	// 		return
	// 	}

	// 	logger.Info("gRPC伺服器啟動", logger.Fields{
	// 		"port": grpcPort,
	// 	})

	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		logger.Error("gRPC伺服器啟動失敗", logger.Fields{
	// 			"error": err.Error(),
	// 		})
	// 	}
	// }()

	// 啟動HTTP伺服器
	go func() {
		logger.Info("HTTP伺服器啟動", logger.Fields{
			"port": cfg.Server.Port,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP伺服器啟動失敗", logger.Fields{
				"error": err.Error(),
			})
			log.Fatal("HTTP伺服器啟動失敗")
		}
	}()

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在關閉伺服器...")

	// 優雅關閉
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 關閉HTTP伺服器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("HTTP伺服器關閉失敗", logger.Fields{
			"error": err.Error(),
		})
	}

	// 關閉gRPC伺服器
	// TODO: 修復 gRPC 服務器
	// grpcServer.GracefulStop()

	// 關閉MQTT連接
	if mqttClient != nil {
		mqttClient.Disconnect()
	}

	logger.Info("伺服器已關閉")
}

// setupMiddlewares 設定中介軟體
func setupMiddlewares(r *gin.Engine, cfg *config.Config, jwtManager *pkgjwt.JWTManager) {
	// 恢復中介軟體
	r.Use(middleware.RecoveryMiddleware())

	// CORS中介軟體
	r.Use(middleware.CORSMiddleware())

	// 日誌中介軟體
	r.Use(middleware.LoggerMiddleware())

	// 開發模式下的額外中介軟體
	if cfg.Environment != "production" {
		r.Use(gin.Logger())
	}
}

// setupRoutes 設定API路由
func setupRoutes(r *gin.Engine, threatIntelHandler *handler.ThreatIntelligenceHandler, collectorHandler *handler.CollectorHandler, authHandler *handler.AuthHandler, hibpHandler *handler.HIBPHandler, jwtManager *pkgjwt.JWTManager) {
	// 健康檢查端點
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"message":   "資安情報平台 API 正常運行",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
			"services": gin.H{
				"http": "running",
				"grpc": "running",
				"mqtt": "available",
			},
		})
	})

	// Swagger文檔
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由群組
	api := r.Group("/api/v1")
	{
		// 認證路由（公開）
		authHandler.RegisterRoutes(api)

		// 需要認證的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuthMiddleware(jwtManager))
		{
			// 威脅情報路由
			threatIntel := authenticated.Group("/threat-intelligence")
			{
				// 基本CRUD操作
				threatIntel.GET("", threatIntelHandler.ListThreats)
				threatIntel.POST("", threatIntelHandler.CreateThreat)
				threatIntel.GET("/:id", threatIntelHandler.GetThreat)
				threatIntel.PUT("/:id", threatIntelHandler.UpdateThreat)
				threatIntel.DELETE("/:id", threatIntelHandler.DeleteThreat)

				// 搜尋和查詢
				threatIntel.GET("/search", threatIntelHandler.SearchThreats)
				threatIntel.GET("/lookup/ip", threatIntelHandler.LookupIP)
				threatIntel.GET("/lookup/domain", threatIntelHandler.LookupDomain)
				
				// 統計和分析
				threatIntel.GET("/statistics", threatIntelHandler.GetStatistics)
				
				// 批量操作
				threatIntel.POST("/batch", threatIntelHandler.BulkCreateThreats)
				threatIntel.PUT("/batch", threatIntelHandler.BulkUpdateThreats)
				threatIntel.DELETE("/batch", threatIntelHandler.BulkDeleteThreats)
			}

			// 收集器路由
			collector := authenticated.Group("/collector")
			{
				collector.POST("/collect-ip", collectorHandler.CollectIPThreatIntel)
				collector.POST("/collect-ips", collectorHandler.CollectBulkIPThreatIntel)
			}

			// HIBP 路由
			hibpHandler.RegisterRoutes(authenticated)
		}
	}
}

// getEnvOrDefault 取得環境變數或預設值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 