package app

import (
	"context"
	"event_booker/internal/config"
	"event_booker/internal/handler"
	"event_booker/internal/middleware"
	"event_booker/internal/repository"
	"event_booker/internal/scheduler"
	"event_booker/internal/usecase"
	"event_booker/migrations"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func Run() {
	server := ginext.New("release")

	server.Use(middleware.ServerMiddleware())

	// Logger
	zlog.InitConsole()

	serviceCfg := config.New()
	if serviceCfg == nil {
		zlog.Logger.Fatal().Msg("Failed to load config")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		serviceCfg.DbUser,
		serviceCfg.DbPassword,
		serviceCfg.DbHost,
		serviceCfg.DbPort,
		serviceCfg.DbName,
	)

	zlog.Logger.Info().Str("dsn", dsn).Msg("Database connection string")

	redisClient := redis.New(serviceCfg.RedisAddr, "", 0)

	repo := repository.New(dsn, &dbpg.Options{
		MaxOpenConns:    serviceCfg.MaxOpenConns,
		MaxIdleConns:    serviceCfg.MaxIdleConns,
		ConnMaxLifetime: serviceCfg.ConnMaxLifetime,
	},
		redisClient,
		serviceCfg,
	)

	shortenerMigrations := migrations.New(repo, serviceCfg)
	if err := shortenerMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	useCase := usecase.New(repo)

	eventHandler := handler.New(useCase)

	schedulerBooker := scheduler.New(useCase)
	schedulerBooker.Start()
	defer schedulerBooker.Stop()

	// API routes
	server.POST("/events", eventHandler.CreateEvent)
	server.GET("/events/all", eventHandler.GetAllEvents)
	server.POST("/events/:id/book", eventHandler.EventBook)
	server.POST("/events/:id/confirm", eventHandler.Confirm)
	server.GET("/events/:id", eventHandler.GetEvent)

	// файлы для фронта
	server.Engine.Static("/web", "./web")
	server.Engine.StaticFile("/", "./web/index.html")

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Error starting server")
	}
}
