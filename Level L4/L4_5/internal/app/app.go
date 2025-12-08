package app

import (
	"api_optimization/internal/handler"
	"api_optimization/internal/usecase"
	"log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Run - страт сервиса
func Run() {
	server := gin.Default()
	pprof.Register(server)

	uc := usecase.New()

	h := handler.New(uc)

	server.POST("/sum", h.Sum)

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
