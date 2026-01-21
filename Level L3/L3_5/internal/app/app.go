package app

import (
	"context"
	"event_booker/internal/config"
	"event_booker/internal/handler"
	"event_booker/internal/repository"
	"event_booker/internal/scheduler"
	"event_booker/internal/usecase"
	"event_booker/migrations"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

// Run - запуск сервиса
func Run() {
	server := ginext.New("release")

	// Logger
	zlog.InitConsole()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	redisClient := redis.New(cfg.RedisAddr, "", 0)

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	},
		redisClient,
	)

	shortenerMigrations := migrations.New(repo, cfg)
	if err := shortenerMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	useCase := usecase.New(repo)

	eventHandler := handler.New(useCase)

	schedulerBooker := scheduler.New(useCase)
	schedulerBooker.Start()
	defer schedulerBooker.Stop()

	zlog.Logger.Info().Msg("Service starting successfully")

	// API routes
	server.POST("/events", eventHandler.CreateEvent)
	server.GET("/events/all", eventHandler.GetAllEvents)
	server.POST("/events/:id/book", eventHandler.EventBook)
	server.POST("/events/:id/confirm", eventHandler.Confirm)
	server.GET("/events/:id", eventHandler.GetEvent)

	// Статические файлы для фронтенда
	server.Engine.Static("/web", "./web")
	server.Engine.StaticFile("/", "./web/index.html")

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Error starting server")
	}
}
