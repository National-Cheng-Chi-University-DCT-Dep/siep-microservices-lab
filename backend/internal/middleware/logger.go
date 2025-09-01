package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// LoggerMiddleware 記錄 HTTP 請求的中介軟體
func LoggerMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 開始時間
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		// 處理請求
		c.Next()

		// 計算耗時
		param := gin.LogFormatterParams{
			Request:    c.Request,
			TimeStamp:  time.Now(),
			Latency:    time.Since(start),
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:   c.Writer.Size(),
			Keys:       c.Keys,
		}

		if raw != "" {
			param.Path = path + "?" + raw
		} else {
			param.Path = path
		}

		// 記錄請求
		fields := pkglogger.Fields{
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"status_code": param.StatusCode,
			"latency_ms":  float64(param.Latency.Nanoseconds()) / 1000000,
			"body_size":   param.BodySize,
			"user_agent":  c.Request.UserAgent(),
		}

		// 如果有錯誤，加入錯誤資訊
		if param.ErrorMessage != "" {
			fields["error"] = param.ErrorMessage
		}

		// 根據狀態碼決定日誌級別
		switch {
		case param.StatusCode >= 500:
			pkglogger.Error("HTTP Request", fields)
		case param.StatusCode >= 400:
			pkglogger.Warn("HTTP Request", fields)
		default:
			pkglogger.Info("HTTP Request", fields)
		}
	})
}

// RequestIDMiddleware 為每個請求生成唯一 ID
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	})
}

// generateRequestID 生成請求 ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成隨機字串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
} 