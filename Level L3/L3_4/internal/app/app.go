package app

import (
	"context"
	"image_processor/internal/config"
	"image_processor/internal/handler"
	appKafka "image_processor/internal/kafka"
	"image_processor/internal/middleware"
	"image_processor/internal/repository"
	"image_processor/internal/usecase"
	"image_processor/migrations"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/zlog"
)

// Run - запуск
func Run() {
	server := ginext.New("release")

	// мидлвара чтобы шли ручки с фронта
	server.Use(middleware.ServerMiddleware())

	zlog.InitConsole()

	// загружаем конфиг
	cfg, err := config.New()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	repo := repository.New(cfg.CreateDsn(), &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}, cfg)

	// миграции
	imgMigrations := migrations.New(repo, cfg)
	if err := imgMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	// usecase
	imgUseCase := usecase.New(repo, cfg)

	// кафка
	consumer := kafka.NewConsumer([]string{cfg.KafkaAddr}, cfg.KafkaTopic, cfg.KafkaGroupID)
	producer := kafka.NewProducer([]string{cfg.KafkaAddr}, cfg.KafkaTopic)

	// создаем очередь
	queue := appKafka.New(consumer, producer, imgUseCase)

	// консюмер
	go func() {
		ctx := context.Background()
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						zlog.Logger.Error().Msgf("Kafka consumer panic: %v", r)
					}
				}()
				zlog.Logger.Info().Msg("Kafka consumer started")
				queue.StartConsumer(ctx)
				zlog.Logger.Warn().Msg("Kafka consumer stopped, restarting in 5s")
			}()
			time.Sleep(5 * time.Second)
		}
	}()

	// инициализация HTTP handler
	imgHandler := handler.New(imgUseCase, queue)

	zlog.Logger.Info().Msg("Service started successfully")

	server.POST("/upload", imgHandler.UploadImage)
	server.POST("/process", imgHandler.PhotoProcessing)
	server.GET("/image/:id/:photo_version", imgHandler.GetProcessedImg)
	server.DELETE("/image/delete", imgHandler.DeletePhoto)

	// для фронта
	server.Static("/static", "./web")
	server.GET("/", func(c *ginext.Context) { c.File("./web/frontend.html") })

	// запуск сервера
	if err := server.Run(":8081"); err != nil {
		zlog.Logger.Fatal().Msg("Failed to run server")
	}
}
