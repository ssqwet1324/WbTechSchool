package main

import (
	"mem_gc_exporter/internal/app"
)

func main() {
	// настройка gc для теста

	//debug.SetGCPercent(100)
	//debug.SetGCPercent(50)
	//debug.SetGCPercent(20)
	//debug.SetGCPercent(-1)

	app.Run()
}
