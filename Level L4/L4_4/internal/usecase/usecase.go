package usecase

import (
	"log/slog"
	"mem_gc_exporter/internal/entity"
	"runtime"
)

// UseCase - бизнес логика
type UseCase struct {
	logger *slog.Logger
}

// New - конструктор
func New(logger *slog.Logger) *UseCase {
	return &UseCase{
		logger: logger,
	}
}

// GetStatistics - получить статистику
func (u *UseCase) GetStatistics() *entity.Metrics {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	metrics := &entity.Metrics{
		Allocations:   mem.Mallocs - mem.Frees,
		GcCount:       mem.NumGC,
		MemoryUsed:    mem.Alloc,
		LastGCPauseNs: mem.PauseNs[(mem.NumGC+255)%256],
	}

	u.logger.Info("GetStatistics", "metrics success", metrics)

	return metrics
}
