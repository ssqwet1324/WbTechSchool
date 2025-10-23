package kafka

import (
	"context"

	"github.com/wb-go/wbf/zlog"
)

// SendMessage отправляет одно сообщение в Kafka
func (q *Queue) SendMessage(ctx context.Context, key, value []byte) error {
	if err := q.producer.Send(ctx, key, value); err != nil {
		zlog.Logger.Error().Err(err).Str("key", string(key)).Msg("PhotoProducer: failed to send message")
		return err
	}

	return nil
}
