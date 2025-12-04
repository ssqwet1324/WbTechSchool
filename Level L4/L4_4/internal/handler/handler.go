package handler

import (
	"log/slog"
	"mem_gc_exporter/internal/usecase"
	"mem_gc_exporter/pkg/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler - обработчик ручек
type Handler struct {
	uc          *usecase.UseCase
	logger      *slog.Logger
	promMetrics *prometheus.Metrics
}

// New - конструктор
func New(uc *usecase.UseCase, logger *slog.Logger, promMetrics *prometheus.Metrics) *Handler {
	return &Handler{
		uc:          uc,
		logger:      logger,
		promMetrics: promMetrics,
	}
}

// Metrics - handler для Prometheus метрик
func (h *Handler) Metrics(ctx *gin.Context) {
	// актуальные метрики
	stats := h.uc.GetStatistics()

	// Обновляем Prometheus метрики
	h.promMetrics.UpdateFromEntity(
		"runtime",
		stats.Allocations,
		stats.GcCount,
		stats.MemoryUsed,
		stats.LastGCPauseNs,
	)

	// Отдаем стандартный Prometheus handler
	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
}
