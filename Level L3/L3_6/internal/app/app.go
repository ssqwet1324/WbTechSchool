package app

import (
	"context"
	"fmt"
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

func Run() {
	server := ginext.New("release")

	server.Use(middleware.ServerMiddleware())

	//логгер
	zlog.InitConsole()

	serviceCfg := config.New()
	if serviceCfg == nil {
		zlog.Logger.Fatal().Msg("Failed to load config")
		return
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		serviceCfg.DbUser,
		serviceCfg.DbPassword,
		serviceCfg.DbHost,
		serviceCfg.DbPort,
		serviceCfg.DbName,
	)

	repo := repository.New(dsn, &dbpg.Options{
		MaxOpenConns:    serviceCfg.MaxOpenConns,
		MaxIdleConns:    serviceCfg.MaxIdleConns,
		ConnMaxLifetime: serviceCfg.ConnMaxLifetime,
	},
	)

	itemsMigrations := migrations.New(repo, serviceCfg)
	if err := itemsMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	useCase := usecase.New(repo)

	itemsHandler := handler.New(useCase)

	server.POST("/items", itemsHandler.CreateItem)
	server.GET("/items", itemsHandler.GetItems)
	server.PUT("/items/:id", itemsHandler.UpdateItem)
	server.DELETE("/items/:id", itemsHandler.DeleteItem)
	server.GET("/analytics", itemsHandler.GetAnalytics)

	// Статические файлы для фронтенда
	server.Engine.Static("/web", "./web")
	server.Engine.StaticFile("/", "./web/index.html")

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to run server")
	}
}
