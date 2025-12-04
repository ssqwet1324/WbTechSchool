package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics - структура метрик
type Metrics struct {
	allocations   *prometheus.GaugeVec
	gcCount       *prometheus.GaugeVec
	memoryUsed    *prometheus.GaugeVec
	lastGCPauseNs *prometheus.GaugeVec
}

// NewMetrics - конструктор
func NewMetrics(name string) *Metrics {
	return &Metrics{
		allocations: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: name + "_allocations",
			Help: "Current number of allocations",
		}, []string{"type"}),
		gcCount: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: name + "_gc_count",
			Help: "Total number of garbage collections",
		}, []string{"type"}),
		memoryUsed: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: name + "_memory_used_bytes",
			Help: "Current memory usage in bytes",
		}, []string{"type"}),
		lastGCPauseNs: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: name + "_last_gc_pause_ns",
			Help: "Duration of the last garbage collection pause in nanoseconds",
		}, []string{"type"}),
	}
}

// Allocate - аллокации
func (m *Metrics) Allocate(label string, value uint64) {
	m.allocations.WithLabelValues(label).Set(float64(value))
}

// GCCount - количество сборок мусора
func (m *Metrics) GCCount(label string, value uint32) {
	m.gcCount.WithLabelValues(label).Set(float64(value))
}

// MemoryUsed - использование памяти
func (m *Metrics) MemoryUsed(label string, bytes uint64) {
	m.memoryUsed.WithLabelValues(label).Set(float64(bytes))
}

// LastGCPauseNs - последнее время GC
func (m *Metrics) LastGCPauseNs(label string, ns uint64) {
	m.lastGCPauseNs.WithLabelValues(label).Set(float64(ns))
}

// UpdateFromEntity - обновить все метрики из entity.Metrics
func (m *Metrics) UpdateFromEntity(label string, allocations uint64, gcCount uint32, memoryUsed uint64, lastGCPauseNs uint64) {
	m.allocations.WithLabelValues(label).Set(float64(allocations))
	m.gcCount.WithLabelValues(label).Set(float64(gcCount))
	m.memoryUsed.WithLabelValues(label).Set(float64(memoryUsed))
	m.lastGCPauseNs.WithLabelValues(label).Set(float64(lastGCPauseNs))
}
