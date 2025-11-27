package middleware

import (
	"L4_3/internal/log"
	"time"

	"github.com/gin-gonic/gin"
)

// LogRequest - логируем запрос и время выполнения
func LogRequest(logger *log.Log) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		if logger != nil {
			logger.AsyncMessagef("[HTTP] start method=%s url=%s time=%s",
				ctx.Request.Method,
				ctx.Request.RequestURI,
				start.Format(time.RFC3339),
			)
		}

		ctx.Next()

		if logger != nil {
			duration := time.Since(start)
			logger.AsyncMessagef("[HTTP] end method=%s url=%s duration=%s status=%d",
				ctx.Request.Method,
				ctx.Request.RequestURI,
				duration.String(),
				ctx.Writer.Status(),
			)
		}
	}
}
