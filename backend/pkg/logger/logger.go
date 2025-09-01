package logger

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel 日誌級別
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String 返回日誌級別的字串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 結構化日誌器
type Logger struct {
	level  LogLevel
	prefix string
	logger *log.Logger
}

// Fields 日誌欄位
type Fields map[string]interface{}

var globalLogger *Logger

// Init 初始化全域日誌器
func Init(level string) {
	logLevel := parseLogLevel(level)
	globalLogger = &Logger{
		level:  logLevel,
		prefix: "[SECURITY-INTEL]",
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// parseLogLevel 解析日誌級別
func parseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// WithFields 添加欄位到日誌
func WithFields(fields Fields) *Logger {
	return &Logger{
		level:  globalLogger.level,
		prefix: globalLogger.prefix,
		logger: globalLogger.logger,
	}
}

// Debug 記錄調試級別日誌
func Debug(msg string, fields ...Fields) {
	if globalLogger != nil {
		globalLogger.log(DEBUG, msg, fields...)
	}
}

// Info 記錄信息級別日誌
func Info(msg string, fields ...Fields) {
	if globalLogger != nil {
		globalLogger.log(INFO, msg, fields...)
	}
}

// Warn 記錄警告級別日誌
func Warn(msg string, fields ...Fields) {
	if globalLogger != nil {
		globalLogger.log(WARN, msg, fields...)
	}
}

// Error 記錄錯誤級別日誌
func Error(msg string, fields ...Fields) {
	if globalLogger != nil {
		globalLogger.log(ERROR, msg, fields...)
	}
}

// Fatal 記錄致命級別日誌並退出程式
func Fatal(msg string, fields ...Fields) {
	if globalLogger != nil {
		globalLogger.log(FATAL, msg, fields...)
		os.Exit(1)
	}
}

// log 內部日誌記錄方法
func (l *Logger) log(level LogLevel, msg string, fields ...Fields) {
	if level < l.level {
		return
	}

	var fieldStr string
	if len(fields) > 0 {
		fieldStr = l.formatFields(fields[0])
	}

	l.logger.Printf("%s [%s] %s %s", l.prefix, level.String(), msg, fieldStr)
}

// formatFields 格式化欄位
func (l *Logger) formatFields(fields Fields) string {
	if len(fields) == 0 {
		return ""
	}

	var parts []string
	for k, v := range fields {
		parts = append(parts, k+"="+l.formatValue(v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// formatValue 格式化值
func (l *Logger) formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return "\"" + val + "\""
	case int, int8, int16, int32, int64:
		return string(rune(val.(int)))
	case float32, float64:
		return string(rune(int(val.(float64))))
	case bool:
		if val {
			return "true"
		}
		return "false"
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return "\"" + string(rune(val.(int))) + "\""
	}
}

// GinLogger 回傳 Gin 框架使用的日誌中介軟體
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 處理請求
		c.Next()

		// 記錄請求日誌
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := Fields{
			"method":      method,
			"path":        path,
			"status":      statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"body_size":   bodySize,
			"user_agent":  c.Request.UserAgent(),
		}

		if statusCode >= 400 {
			Error("HTTP Request Error", fields)
		} else {
			Info("HTTP Request", fields)
		}
	}
}

// GinRecovery 回傳 Gin 框架使用的恢復中介軟體
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fields := Fields{
					"error":     err,
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"client_ip": c.ClientIP(),
				}
				Error("Panic Recovery", fields)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// FromContext 從 context 中提取請求 ID
func FromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// GetLogger 取得全域日誌器
func GetLogger() *Logger {
	return globalLogger
} 