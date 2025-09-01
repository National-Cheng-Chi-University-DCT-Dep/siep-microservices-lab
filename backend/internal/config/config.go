package config

import (
	"os"
	"strconv"
)

// Config 應用程式配置結構
type Config struct {
	Environment string         `json:"environment"`
	Server      ServerConfig   `json:"server"`
	Database    DatabaseConfig `json:"database"`
	LogLevel    string         `json:"log_level"`
	JWT         JWTConfig      `json:"jwt"`
	External    ExternalConfig `json:"external"`
}

// ServerConfig 伺服器配置
type ServerConfig struct {
	Port    string `json:"port"`
	Host    string `json:"host"`
	Timeout int    `json:"timeout"`
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	DSN             string `json:"dsn"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `json:"secret"`
	Expiration int    `json:"expiration"`
}

// ExternalConfig 外部服務配置
type ExternalConfig struct {
	AbuseIPDBKey string `json:"abuse_ipdb_key"`
	HIBPAPIKey   string `json:"hibp_api_key"`
	StripeKey    string `json:"stripe_key"`
	LLMAPIKey    string `json:"llm_api_key"`
}

// Load 載入配置
func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			Host:    getEnv("SERVER_HOST", "localhost"),
			Timeout: getEnvAsInt("SERVER_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", "postgres://user:password@localhost/security_intel?sslmode=disable"),
			MaxOpenConns:    getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvAsInt("DATABASE_CONN_MAX_LIFETIME", 5),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getEnvAsInt("JWT_EXPIRATION", 24), // 24 小時
		},
		External: ExternalConfig{
			AbuseIPDBKey: getEnv("ABUSE_IPDB_KEY", ""),
			HIBPAPIKey:   getEnv("HIBP_API_KEY", ""),
			StripeKey:    getEnv("STRIPE_KEY", ""),
			LLMAPIKey:    getEnv("LLM_API_KEY", ""),
		},
	}

	return cfg, nil
}

// getEnv 取得環境變數，如果不存在則回傳預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 取得環境變數作為整數，如果不存在則回傳預設值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 取得環境變數作為布林值，如果不存在則回傳預設值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
} 