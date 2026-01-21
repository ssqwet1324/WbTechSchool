package app

import (
	"context"
	"shortener/internal/cache"
	service "shortener/internal/config"
	"shortener/internal/handler"
	"shortener/internal/repository"
	"shortener/internal/usecase"
	"shortener/migrations"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

// Run - запуск сервиса
func Run() {
	server := ginext.New()

	zlog.Init()

	cfg, err := service.New()
	if err != nil {
		panic("error initializing service" + err.Error())
	}

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	})

	shortenerMigrations := migrations.New(repo, cfg)
	if err := shortenerMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	redisClient := redis.New(cfg.RedisAddr, "", 0)
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Error closing redis connection")
		}
	}(redisClient)

	shortenerCache := cache.New(*redisClient)

	shortenerUseCase := usecase.New(repo, shortenerCache)

	shortenerHandler := handler.New(shortenerUseCase)

	server.POST("/shorten", shortenerHandler.CreateShorten)
	server.GET("/s/:short_url", shortenerHandler.RedirectToShorten)
	server.GET("/analytics/:short_url", shortenerHandler.GetAnalytics)

	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("error starting server")
	}
}
