package middleware

import (
	"net/http"
	"strings"
	"warehouse_control/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wb-go/wbf/ginext"
)

// ServerMiddleware - мидлвара для ролей
func ServerMiddleware(cfg *config.ServiceConfig) ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		// берем заголовок
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			return
		}

		// делим заголовок 2 на части
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			return
		}

		// берем токен
		tokenStr := parts[1]

		// парсим токен
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// читаем роль в jwt
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role, ok := claims["role"]
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				return
			}
			switch role {
			case "admin":
				// доступ к любой ручке
				ctx.Next()
				return
			case "manager":
				// изменение и просмотр
				if ctx.Request.Method == http.MethodGet ||
					ctx.Request.Method == http.MethodPut ||
					ctx.Request.Method == http.MethodPatch {
					ctx.Next()
					return
				}
			case "viewer":
				// только GET, Только просмотр
				if ctx.Request.Method == http.MethodGet {
					ctx.Next()
					return
				}
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "viewer: read-only access"})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
				return
			}
		}
	}
}

// CorsMiddleware - мидлвара чтобы работали запросы
func CorsMiddleware() ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
