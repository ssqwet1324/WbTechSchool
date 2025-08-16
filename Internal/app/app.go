package app

import (
	"WbDemoProject/Internal/config"
	"WbDemoProject/Internal/handler"
	"WbDemoProject/Internal/kafka"
	"WbDemoProject/Internal/migrations"
	"WbDemoProject/Internal/repository"
	"WbDemoProject/Internal/usecase"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	server := gin.Default()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	migration := migrations.New(repo)
	err = migration.InitTables(context.Background())
	if err != nil {
		log.Fatal("database dont created", err)
	}

	usecase := usecase.New(repo)

	ordersConsumer := kafka.New([]string{"kafka:9092"}, "orders", "1")
	defer func(ordersConsumer *kafka.Consumer) {
		err := ordersConsumer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ordersConsumer)
	go func() {
		for {
			err := ordersConsumer.StartConsumer(context.Background(), usecase)
			if err != nil {
				log.Printf("Kafka not ready, retrying in 5s: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}()

	orderHandler := handler.New(usecase)

	server.GET("/order/:order_uid", orderHandler.GetOrder)

	if err := server.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
