package app

import (
	"fmt"
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

func Run() {
	server := ginext.New("release")

	server.Use(middleware.CorsMiddleware())

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

	allMigrations := migrations.New(repo, serviceCfg)
	if err := allMigrations.RunMigrations(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	useCase := usecase.New(repo, serviceCfg)

	productHandler := handler.New(useCase)

	// публичные ручки
	userGroup := server.Group("/user")
	userGroup.POST("/create", productHandler.CreateUser)
	userGroup.POST("/login", productHandler.Login)

	//ручки по ролям
	roleGroup := server.Group("/api/v1", middleware.ServerMiddleware(serviceCfg))
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
