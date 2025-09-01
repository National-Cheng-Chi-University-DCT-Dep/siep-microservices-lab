package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 配置 CORS 政策
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 從環境變數取得允許的來源，預設為 localhost
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000"
		}

		origin := c.Request.Header.Get("Origin")
		
		// 檢查來源是否被允許
		if isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// isOriginAllowed 檢查來源是否在允許清單中
func isOriginAllowed(origin, allowedOrigins string) bool {
	if origin == "" {
		return false
	}
	
	origins := strings.Split(allowedOrigins, ",")
	for _, allowedOrigin := range origins {
		if strings.TrimSpace(allowedOrigin) == origin {
			return true
		}
	}
	return false
} 