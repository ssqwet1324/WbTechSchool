package middleware

import (
	"github.com/wb-go/wbf/ginext"
)

// ServerMiddleware - мидлвара чтобы работали запросы
func ServerMiddleware() ginext.HandlerFunc {
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
