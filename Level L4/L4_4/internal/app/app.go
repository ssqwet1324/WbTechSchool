package app

import (
	"log"
	"log/slog"
	"mem_gc_exporter/internal/handler"
	"mem_gc_exporter/internal/middleware"
	"mem_gc_exporter/internal/usecase"
	"mem_gc_exporter/pkg/prometheus"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Run() {
	server := gin.Default()

	// добавляем pprof
	pprof.Register(server)

	// Инициализация логгера
	logger := slog.Default()

	uc := usecase.New(logger)

	// Инициализация Prometheus метрик
	promMetrics := prometheus.NewMetrics("mem_gc_exporter")

	metricsHandler := handler.New(uc, logger, promMetrics)

	// Подключение middleware
	server.Use(middleware.PrometheusMiddleware())

	// ручка для метрик
	server.GET("/metrics", metricsHandler.Metrics)

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
