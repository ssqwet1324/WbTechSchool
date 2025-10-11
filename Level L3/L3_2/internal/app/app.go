package app

import (
	"context"
	"fmt"
	"log"
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

func Run() {
	server := ginext.New()

	// Logger
	zlog.Init()

	serviceCfg := service.New()
	if serviceCfg == nil {
		zlog.Logger.Fatal().Msg("Failed to load config")
		return
	}

	// Логируем конфигурацию для отладки
	zlog.Logger.Info().
		Str("db_host", serviceCfg.DbHost).
		Int("db_port", serviceCfg.DbPort).
		Str("db_name", serviceCfg.DbName).
		Str("db_user", serviceCfg.DbUser).
		Msg("Loaded configuration")

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

	shortenerMigrations := migrations.New(repo, serviceCfg)
	if err := shortenerMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	redisClient := redis.New(serviceCfg.RedisAddr, "", 0)
	shortenerCache := cache.New(*redisClient)

	shortenerUseCase := usecase.New(repo, shortenerCache)

	shortenerHandler := handler.New(shortenerUseCase)

	server.POST("/shorten", shortenerHandler.CreateShorten)
	server.GET("/s/:short_url", shortenerHandler.RedirectToShorten)
	server.GET("/analytics/:short_url", shortenerHandler.GetAnalytics)

	if err := server.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
