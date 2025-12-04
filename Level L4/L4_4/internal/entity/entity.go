package entity

// Metrics - хранит основные метрики памяти и сборщика мусора
type Metrics struct {
	Allocations   uint64 `json:"allocations"`   // количество аллокаций
	GcCount       uint32 `json:"gc_count"`      // количество сборок мусора
	MemoryUsed    uint64 `json:"memory_used"`   // используемая память
	LastGCPauseNs uint64 `json:"last_gc_pause"` // последнее время GC
}
