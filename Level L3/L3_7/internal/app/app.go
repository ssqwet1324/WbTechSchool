package app

import (
	"warehouse_control/internal/config"
	"warehouse_control/internal/handler"
	"warehouse_control/internal/middleware"
	"warehouse_control/internal/repository"
	"warehouse_control/internal/usecase"
	"warehouse_control/migrations"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// Run - запуск сервиса
func Run() {
	server := ginext.New("release")

	server.Use(middleware.CorsMiddleware())

	zlog.InitConsole()

	cfg, err := config.New()
	if cfg == nil || err != nil {
		panic("Error loading config")
	}

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	},
	)

	allMigrations := migrations.New(repo, cfg)
	if err := allMigrations.RunMigrations(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	useCase := usecase.New(repo, cfg)

	productHandler := handler.New(useCase)

	zlog.Logger.Info().Msg("Service started successfully")

	// публичные ручки
	userGroup := server.Group("/user")
	userGroup.POST("/create", productHandler.CreateUser)
	userGroup.POST("/login", productHandler.Login)

	//ручки по ролям
	roleGroup := server.Group("/api/v1", middleware.ServerMiddleware(cfg))
	roleGroup.POST("/items", productHandler.CreateProduct)
	roleGroup.GET("/items/:product_name", productHandler.GetProduct)
	roleGroup.GET("/items", productHandler.GetAllProduct)
	roleGroup.PUT("/items", productHandler.UpdateProduct)
	roleGroup.DELETE("/items/:product_name", productHandler.DeleteProduct)
	roleGroup.GET("/product/logs/:product_id", productHandler.GetLogsByProductID)
	roleGroup.GET("/product/save/:product_id", productHandler.SaveProductLogToCSV)

	// Файлы для фронта
	server.Engine.Static("/css", "./web/css")
	server.Engine.Static("/js", "./web/js")
	server.Engine.StaticFile("/", "./web/index.html")

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Error starting server")
	}
}
