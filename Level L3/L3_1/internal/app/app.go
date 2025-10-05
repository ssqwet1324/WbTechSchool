package app

import (
	"L3_1/internal/cache"
	"L3_1/internal/config"
	"L3_1/internal/handler"
	"L3_1/internal/migrations"
	"L3_1/internal/rabbit"
	"L3_1/internal/repository"
	"L3_1/internal/usecase"
	"L3_1/internal/worker"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

// Run запускает HTTP-сервис уведомлений и фонового воркера
func Run() {
	server := ginext.New()

	zlog.Init()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		strconv.Itoa(cfg.DbPort),
		cfg.DbName,
	)

	repo := repository.New(dsn, &dbpg.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
	})

	// --- Миграции ---
	migration := migrations.New(repo, cfg)
	if err := migration.InitNotifyTable(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	// --- Redis ---
	redisAddr := "redis:6379"
	if cfg.RedisAddr != "" {
		redisAddr = cfg.RedisAddr
	}
	redisClient := redis.New(redisAddr, "", 0)
	cacheNotify := cache.New(redisClient)

	useCaseNotify := usecase.New(repo, cacheNotify)

	// --- RabbitMQ ---
	rabbitURL := "amqp://guest:guest@rabbitmq:5672/"
	if cfg.RabbitURL != "" {
		rabbitURL = cfg.RabbitURL
	}

	conn, err := rabbitmq.Connect(rabbitURL, cfg.MaxRetries, 2*time.Second)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось подключиться к RabbitMQ")
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось открыть канал RabbitMQ")
	}
	defer channel.Close()

	rmqPublisher := rabbitmq.NewPublisher(channel, "")
	publisher := rabbit.NewPublisher(rmqPublisher)

	w := worker.New(context.Background(), useCaseNotify, publisher)
	go w.Run()

	handlerNotify := handler.New(useCaseNotify, w)
	server.POST("/notify", handlerNotify.CreateNotification)
	server.DELETE("/notify/:notifyID", handlerNotify.DeleteNotification)
	server.GET("/notify/:notifyID", handlerNotify.CheckStatusNotification)
	server.GET("/notifications/:userID", handlerNotify.GetAllNotifications)

	// --- HTTP Server ---
	if err := server.Run(":8081"); err != nil {
		panic(err)
	}
}
