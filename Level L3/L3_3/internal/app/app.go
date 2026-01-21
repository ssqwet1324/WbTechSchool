package app

import (
	"comment_tree/internal/config"
	"comment_tree/internal/handler"
	"comment_tree/internal/middleware"
	"comment_tree/internal/repository"
	"comment_tree/internal/usecase"
	"comment_tree/migrations"
	"context"

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

	cfg, err := config.New()
	if err != nil {
		panic("Error loading config" + err.Error())
	}

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	})

	commentMigrations := migrations.New(repo, cfg)
	if err := commentMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	commentUseCase := usecase.New(repo)
	commentHandler := handler.New(commentUseCase)

	zlog.Logger.Info().Msg("Service started successfully")

	server.POST("/comments", commentHandler.AddComment)
	server.GET("/comments", commentHandler.GetComments)
	server.DELETE("/comments/:comment_id", commentHandler.DeleteComment)

	server.POST("/comments/search", commentHandler.SearchComment)
	server.GET("/comments/parents", commentHandler.GetParentComments)

	// корневая страница для фронта
	server.Static("/static", "./web")
	server.GET("/", func(c *ginext.Context) { c.File("./web/index.html") })

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to run server")
	}
}
