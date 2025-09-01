package middleware

import (
	"github.com/gin-gonic/gin"
)

// SetupMiddlewares 設定所有中介軟體
func SetupMiddlewares(r *gin.Engine) {
	// 順序很重要：Recovery -> RequestID -> Logger -> CORS
	r.Use(RecoveryMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware())
}

// SetupProductionMiddlewares 設定生產環境專用的中介軟體
func SetupProductionMiddlewares(r *gin.Engine) {
	// 生產環境關閉 Gin 的預設日誌
	gin.SetMode(gin.ReleaseMode)
	
	// 設定基本中介軟體
	SetupMiddlewares(r)
	
	// 這裡可以加入生產環境專用的中介軟體
	// 例如：限流、認證、監控等
}

// SetupDevelopmentMiddlewares 設定開發環境專用的中介軟體
func SetupDevelopmentMiddlewares(r *gin.Engine) {
	// 開發環境保持 Gin 的詳細日誌
	gin.SetMode(gin.DebugMode)
	
	// 設定基本中介軟體
	SetupMiddlewares(r)
	
	// 這裡可以加入開發環境專用的中介軟體
	// 例如：詳細的除錯資訊、熱重載等
} 