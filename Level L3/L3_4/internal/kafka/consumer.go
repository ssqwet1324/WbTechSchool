package kafka

import (
	"context"
	"encoding/json"
	"image_processor/internal/entity"
	"image_processor/internal/usecase"

	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/zlog"
)

type Queue struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
	uc       *usecase.UseCase
}

func New(c *kafka.Consumer, pr *kafka.Producer, uc *usecase.UseCase) *Queue {
	return &Queue{
		consumer: c,
		producer: pr,
		uc:       uc,
	}
}

// StartConsumer - работа консюмера
func (q *Queue) StartConsumer(ctx context.Context) {
	for {
		// Получаем сообщение
		msg, err := q.consumer.Fetch(ctx)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("StartConsumer: failed to fetch message")
			break
		}

		var photo entity.PhotoInfo
		if err := json.Unmarshal(msg.Value, &photo); err != nil {
			zlog.Logger.Error().Err(err).Msg("StartConsumer: failed unmarshalling")
			if commitErr := q.consumer.Commit(ctx, msg); commitErr != nil {
				zlog.Logger.Error().Err(commitErr).Msg("StartConsumer: failed to commit message after unmarshal error")
			}
			continue
		}

		// Обрабатываем фото
		if _, err := q.uc.PhotoProcessing(ctx, photo); err != nil {
			zlog.Logger.Error().Err(err).Str("photoID", photo.PhotoID).Msg("StartConsumer: failed processing photo")
			if commitErr := q.consumer.Commit(ctx, msg); commitErr != nil {
				zlog.Logger.Error().Err(commitErr).Msg("StartConsumer: failed to commit message after processing error")
			}
			continue
		}

		// Успешная обработка
		if err := q.consumer.Commit(ctx, msg); err != nil {
			zlog.Logger.Error().Err(err).Msg("StartConsumer: failed to commit message")
		}
	}
}
