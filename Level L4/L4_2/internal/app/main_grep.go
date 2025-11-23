package app

import (
	"bufio"

	"L4.2/internal/entity"
	"L4.2/internal/handler"
	"L4.2/internal/service"
	"github.com/gin-gonic/gin"
)

// RunWorker - Запускаем воркера
func RunWorker(flag entity.ServerFlags) error {
	server := gin.Default()
	h := handler.NewWorkerHandler()
	server.POST("/search", h.HandleSearch)
	return server.Run(flag.Addr)
}

// RunServer - Запускаем сервис лидера
func RunServer(flag entity.ServerFlags, options entity.Options, pattern string, scanner *bufio.Scanner) error {
	// Читаем все строки
	lines, err := service.ReadAll(scanner)
	if err != nil {
		return err
	}

	// Парсим адреса воркеров
	peers := service.ParsePeers(flag.Peers)

	// Выполняем распределённый поиск
	return service.Execute(lines, peers, options, pattern, flag.Quorum)
}
