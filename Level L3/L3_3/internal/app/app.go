package app

import (
	"comment_tree/internal/config"
	"comment_tree/internal/handler"
	"comment_tree/internal/middleware"
	"comment_tree/internal/repository"
	"comment_tree/internal/usecase"
	"comment_tree/migrations"
	"context"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// Run - запуск сервиса
func Run() {
	server := ginext.New("release")

	// мидлвара чтобы шли ручки с фронта
	server.Use(middleware.ServerMiddleware())

	// запускаем логгер
	zlog.InitConsole()

	serviceCfg := config.New()
	if serviceCfg == nil {
		zlog.Logger.Fatal().Msg("Failed to load config")
		return
	}

	// че в логах
	zlog.Logger.Info().
		Str("db_host", serviceCfg.DbHost).
		Int("db_port", serviceCfg.DbPort).
		Str("db_name", serviceCfg.DbName).
		Str("db_user", serviceCfg.DbUser).
		Msg("Loaded configuration")

	// подключение к бд
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		serviceCfg.DbUser,
		serviceCfg.DbPassword,
		serviceCfg.DbHost,
		serviceCfg.DbPort,
		serviceCfg.DbName,
	)

	zlog.Logger.Info().Str("dsn", dsn).Msg("Database connection string")

	repo := repository.New(dsn, &dbpg.Options{
		MaxOpenConns:    serviceCfg.MaxOpenConns,
		MaxIdleConns:    serviceCfg.MaxIdleConns,
		ConnMaxLifetime: serviceCfg.ConnMaxLifetime,
	})

	commentMigrations := migrations.New(repo, serviceCfg)
	if err := commentMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	commentUseCase := usecase.New(repo)

	commentHandler := handler.New(commentUseCase)

	server.POST("/comments", commentHandler.AddComment)
	server.GET("/comments", commentHandler.GetComments)
	server.DELETE("/comments/:comment_id", commentHandler.DeleteComment)

	server.POST("/comments/search", commentHandler.SearchComment)
	server.GET("/comments/parents", commentHandler.GetParentComments)

	// Раздача статики и корневая страница для фронта
	server.Static("/static", "./web")
	server.GET("/", func(c *ginext.Context) { c.File("./web/index.html") })

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Msg("Failed to run server")
	}
}
