package app

import (
	"context"
	"sales_tracker/internal/config"
	"sales_tracker/internal/handler"
	"sales_tracker/internal/middleware"
	"sales_tracker/internal/repository"
	"sales_tracker/internal/usecase"
	"sales_tracker/migrations"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// Run - запуск сервиса
func Run() {
	server := ginext.New("release")

	server.Use(middleware.ServerMiddleware())

	zlog.InitConsole()

	cfg, err := config.New()
	if cfg == nil || err != nil {
		zlog.Logger.Fatal().Msg("Failed to load config")
		return
	}

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	},
	)

	itemsMigrations := migrations.New(repo, cfg)
	if err := itemsMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	useCase := usecase.New(repo)

	itemsHandler := handler.New(useCase)

	zlog.Logger.Info().Msg("Service started successfully")

	server.POST("/items", itemsHandler.CreateItem)
	server.GET("/items", itemsHandler.GetItems)
	server.PUT("/items/:id", itemsHandler.UpdateItem)
	server.DELETE("/items/:id", itemsHandler.DeleteItem)
	server.GET("/analytics", itemsHandler.GetAnalytics)
	server.POST("/analytics/csv", itemsHandler.SaveAnalyticsToCSV)

	// Статические файлы для фронтенда
	server.Engine.Static("/web", "./web")
	server.Engine.StaticFile("/", "./web/index.html")

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to run server")
	}
}
