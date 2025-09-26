package app

import (
	"L2_18/internal/handler"
	"L2_18/internal/middleware"
	"L2_18/internal/repository"
	"L2_18/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Run - старт сервиса
func Run() {
	service := gin.Default()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	service.Use(middleware.LogRequest(logger))

	repo := repository.New(logger)

	usecase := usecase.New(repo, logger)

	eventHandler := handler.New(usecase, logger)

	service.GET("/events_for_day/:user_id/:date", eventHandler.EventsForDay)
	service.GET("/events_for_week/:user_id/:date", eventHandler.EventsForWeek)
	service.GET("/events_for_month/:user_id/:date", eventHandler.EventsForMonth)

	service.POST("/create_event", eventHandler.CreateEvent)
	service.POST("/update_event", eventHandler.UpdateEvent)
	service.POST("/delete_event", eventHandler.DeleteEvent)

	if err := service.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
