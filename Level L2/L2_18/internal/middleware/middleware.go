package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogRequest - логируем запрос и время выполнения
func LogRequest(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		logger.Info("start request",
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.RequestURI),
			zap.Time("time", start),
		)
		ctx.Next()
		duration := time.Since(start)
		logger.Info("end request",
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.RequestURI),
			zap.Duration("duration", duration),
		)
	}
}
