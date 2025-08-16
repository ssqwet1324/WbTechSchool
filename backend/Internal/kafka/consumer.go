package kafka

import (
	"WbDemoProject/Internal/entity"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type OrderHandler interface {
	HandleOrder(ctx context.Context, order *entity.Order) error
}

type Consumer struct {
	reader *kafka.Reader
}

func New(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &Consumer{reader: reader}
}

// StartConsumer - читаем сообщения и отправляем в бизнес логику
func (consumer *Consumer) StartConsumer(ctx context.Context, handler OrderHandler) error {
	for {
		msg, err := consumer.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var order entity.Order
		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)

			if err := consumer.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing invalid message: %v", err)
			}
			continue
		}

		// передаем заказ в обработчик
		if err := handler.HandleOrder(ctx, &order); err != nil {
			log.Printf("Error handling order: %v", err)

			continue
		}

		// фиксируем сообщение
		if err := consumer.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}

func (consumer *Consumer) Close() error {
	return consumer.reader.Close()
}
