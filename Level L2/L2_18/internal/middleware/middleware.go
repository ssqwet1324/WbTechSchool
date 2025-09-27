package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger - тут создаем цветной logger
func NewZapLogger() *zap.Logger {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

	return zap.New(core)
}

// LogRequest - логика logger-а
func LogRequest(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.RequestURI
		size := ctx.Writer.Size()

		msg := fmt.Sprintf("%s %s %s | %v | ip=%s | size=%d",
			method, path, colorStatus(statusCode), duration, clientIP, size,
		)

		switch {
		case statusCode >= 500:
			logger.Error(msg)
		case statusCode >= 400:
			logger.Warn(msg)
		default:
			logger.Info(msg)
		}
	}
}

// colorStatus возвращает цветной код статуса
func colorStatus(code int) string {
	switch code / 100 {
	case 2:
		return fmt.Sprintf("\033[32m%d\033[0m", code) // зелёный
	case 3:
		return fmt.Sprintf("\033[36m%d\033[0m", code) // голубой
	case 4:
		return fmt.Sprintf("\033[33m%d\033[0m", code) // жёлтый
	case 5:
		return fmt.Sprintf("\033[31m%d\033[0m", code) // красный
	default:
		return fmt.Sprintf("%d", code)
	}
}
