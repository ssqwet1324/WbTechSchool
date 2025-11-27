package app

import (
	"L4_3/internal/entity"
	"L4_3/internal/handler"
	asynclog "L4_3/internal/log"
	"L4_3/internal/middleware"
	"L4_3/internal/repository"
	"L4_3/internal/usecase"
	"L4_3/internal/worker"
	stdlog "log"
	"time"

	"github.com/gin-gonic/gin"
)

// Run - старт сервиса
func Run() {
	service := gin.Default()

	// единый асинхронный логгер
	logCh := make(chan string, 200)
	appLogger := asynclog.New(logCh)
	appLogger.StartLogger()

	service.Use(middleware.LogRequest(appLogger))

	repo := repository.New(appLogger)

	// общий канал уведомлений
	reminderCh := make(chan entity.Event, 100)

	usecase := usecase.New(repo, reminderCh)

	worker := worker.New(appLogger, usecase, reminderCh)

	// запускаем воркеров
	worker.RunReminderWorker()
	worker.RunCleanerWorker(10 * time.Minute)

	eventHandler := handler.New(usecase, appLogger)

	service.GET("/events_for_day/:user_id/:date", eventHandler.EventsForDay)
	service.GET("/events_for_week/:user_id/:date", eventHandler.EventsForWeek)
	service.GET("/events_for_month/:user_id/:date", eventHandler.EventsForMonth)

	service.POST("/create_event", eventHandler.CreateEvent)
	service.POST("/update_event", eventHandler.UpdateEvent)
	service.POST("/delete_event", eventHandler.DeleteEvent)

	if err := service.Run(":8080"); err != nil {
		stdlog.Fatal(err)
	}
}
