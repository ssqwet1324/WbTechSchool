package app

import (
	"context"
	"fmt"
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

func Run() {
	server := ginext.New("release")

	// мидлвара чтобы шли ручки с фронта
	server.Use(middleware.ServerMiddleware())

	zlog.InitConsole()

	// загружаем конфиг
	serviceCfg := config.New()
	if serviceCfg == nil {
		zlog.Logger.Fatal().Msg("Failed to load config")
		return
	}

	// логируем конфигурацию
	zlog.Logger.Info().
		Str("db_host", serviceCfg.DbHost).
		Int("db_port", serviceCfg.DbPort).
		Str("db_name", serviceCfg.DbName).
		Str("db_user", serviceCfg.DbUser).
		Msg("Loaded configuration")

	// подключение к БД
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
	}, serviceCfg)

	// миграции
	imgMigrations := migrations.New(repo, serviceCfg)
	if err := imgMigrations.InitTables(context.Background()); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Не удалось создать таблицы")
	}

	// usecase
	imgUseCase := usecase.New(repo, serviceCfg)

	// кафка
	//TODO вынести в кфг
	consumer := kafka.NewConsumer([]string{serviceCfg.KafkaAddr}, serviceCfg.KafkaTopic, serviceCfg.KafkaGroupId)
	producer := kafka.NewProducer([]string{serviceCfg.KafkaAddr}, serviceCfg.KafkaTopic)

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
