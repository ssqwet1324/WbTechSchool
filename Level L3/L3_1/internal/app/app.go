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
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

// Run - запускаем сервис
func Run() {
	server := ginext.New()

	// Logger
	zlog.Init()

	// Config
	cfg, err := config.New()
	if err != nil {
		panic("config init error: " + err.Error())
	}

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime * time.Minute,
	})

	// Migrations
	migration := migrations.New(repo, cfg)
	if err := migration.InitNotifyTable(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	// Redis
	redisClient := redis.New(cfg.RedisAddr, "", 0)
	cacheNotify := cache.New(redisClient)

	useCaseNotify := usecase.New(repo, cacheNotify)

	// RabbitMQ
	rabbitURL := cfg.RabbitURL

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

	// Worker
	w := worker.New(context.Background(), useCaseNotify, publisher)
	go w.Run()

	// Handler
	handlerNotify := handler.New(useCaseNotify, w)

	// Handlers
	server.POST("/notify", handlerNotify.CreateNotification)
	server.DELETE("/notify/:notifyID", handlerNotify.DeleteNotification)
	server.GET("/notify/:notifyID", handlerNotify.CheckStatusNotification)
	server.GET("/notifications/:userID", handlerNotify.GetAllNotifications)

	if err := server.Run(":8081"); err != nil {
		panic(err)
	}
}
