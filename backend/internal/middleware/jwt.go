package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/internal/vo"
	pkgjwt "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/jwt"
	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// JWTAuthMiddleware JWT認證中介軟體
func JWTAuthMiddleware(jwtManager *pkgjwt.JWTManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 從Authorization標頭提取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondUnauthorized(c, "MISSING_AUTH_HEADER", "Authorization header is required")
			return
		}

		// 提取令牌
		tokenString, err := pkgjwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			respondUnauthorized(c, "INVALID_AUTH_HEADER", err.Error())
			return
		}

		// 驗證令牌
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			respondUnauthorized(c, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		// 將使用者資訊設定到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token_claims", claims)

		// 記錄認證成功
		pkglogger.Debug("JWT authentication successful", pkglogger.Fields{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		})

		c.Next()
	})
}

// OptionalJWTAuthMiddleware 可選的JWT認證中介軟體
// 如果有提供令牌則驗證，沒有則繼續執行
func OptionalJWTAuthMiddleware(jwtManager *pkgjwt.JWTManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString, err := pkgjwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 設定使用者資訊
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token_claims", claims)

		c.Next()
	})
}

// RequireRoleMiddleware 角色授權中介軟體
func RequireRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			respondForbidden(c, "MISSING_ROLE", "User role not found in context")
			return
		}

		role, ok := userRole.(string)
		if !ok {
			respondForbidden(c, "INVALID_ROLE", "Invalid user role format")
			return
		}

		// 檢查是否有允許的角色
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		respondForbidden(c, "INSUFFICIENT_PERMISSIONS", "Insufficient permissions for this action")
	})
}

// RequireAdminMiddleware 要求管理員權限的中介軟體
func RequireAdminMiddleware() gin.HandlerFunc {
	return RequireRoleMiddleware("admin")
}

// RequirePremiumMiddleware 要求進階用戶權限的中介軟體
func RequirePremiumMiddleware() gin.HandlerFunc {
	return RequireRoleMiddleware("admin", "premium")
}

// APIKeyAuthMiddleware API金鑰認證中介軟體
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 從標頭或查詢參數提取API金鑰
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			respondUnauthorized(c, "MISSING_API_KEY", "API key is required")
			return
		}

		// TODO: 實作API金鑰驗證邏輯
		// 這裡需要查詢資料庫驗證API金鑰的有效性
		
		pkglogger.Debug("API key authentication", pkglogger.Fields{
			"api_key_prefix": maskAPIKey(apiKey),
		})

		c.Next()
	})
}

// CombinedAuthMiddleware 組合認證中介軟體（JWT或API金鑰）
func CombinedAuthMiddleware(jwtManager *pkgjwt.JWTManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 優先檢查JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString, err := pkgjwt.ExtractTokenFromHeader(authHeader)
			if err == nil {
				claims, err := jwtManager.VerifyToken(tokenString)
				if err == nil {
					// JWT認證成功
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("email", claims.Email)
					c.Set("role", claims.Role)
					c.Set("auth_method", "jwt")
					c.Next()
					return
				}
			}
		}

		// 檢查API金鑰
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey != "" {
			// TODO: 驗證API金鑰
			c.Set("auth_method", "api_key")
			c.Next()
			return
		}

		respondUnauthorized(c, "MISSING_AUTH", "JWT token or API key is required")
	})
}

// respondUnauthorized 回應未授權錯誤
func respondUnauthorized(c *gin.Context, code string, message string) {
	errorVO := vo.ErrorVO{
		Code:    code,
		Message: message,
	}

	response := vo.BaseResponse{
		Success: false,
		Message: "Authentication required",
		Error:   &errorVO,
	}

	c.JSON(http.StatusUnauthorized, response)
	c.Abort()
}

// respondForbidden 回應禁止存取錯誤
func respondForbidden(c *gin.Context, code string, message string) {
	errorVO := vo.ErrorVO{
		Code:    code,
		Message: message,
	}

	response := vo.BaseResponse{
		Success: false,
		Message: "Access forbidden",
		Error:   &errorVO,
	}

	c.JSON(http.StatusForbidden, response)
	c.Abort()
}

// maskAPIKey 遮蔽API金鑰敏感資訊
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
} 