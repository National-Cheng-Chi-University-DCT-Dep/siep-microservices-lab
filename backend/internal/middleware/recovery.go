package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// RecoveryMiddleware 恢復 panic 的中介軟體
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 取得堆疊追蹤
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				stackTrace := string(stack[:length])

				// 記錄 panic 詳細資訊
				fields := pkglogger.Fields{
					"error":       fmt.Sprintf("%v", err),
					"method":      c.Request.Method,
					"path":        c.Request.URL.Path,
					"client_ip":   c.ClientIP(),
					"user_agent":  c.Request.UserAgent(),
					"stack_trace": cleanStackTrace(stackTrace),
				}

				// 如果有請求 ID，也記錄
				if requestID, exists := c.Get("request_id"); exists {
					fields["request_id"] = requestID
				}

				pkglogger.Error("Panic recovered", fields)

				// 回傳 500 錯誤
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred",
				})
				c.Abort()
			}
		}()

		c.Next()
	})
}

// cleanStackTrace 清理堆疊追蹤，移除不必要的資訊
func cleanStackTrace(stack string) string {
	lines := strings.Split(stack, "\n")
	var cleanedLines []string
	
	for _, line := range lines {
		// 跳過 runtime 和 gin 內部的行
		if strings.Contains(line, "runtime/") || 
		   strings.Contains(line, "gin-gonic/gin") ||
		   strings.Contains(line, "net/http") {
			continue
		}
		
		cleanedLines = append(cleanedLines, line)
		
		// 限制堆疊追蹤的長度
		if len(cleanedLines) > 20 {
			break
		}
	}
	
	return strings.Join(cleanedLines, "\n")
} 