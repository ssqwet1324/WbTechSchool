package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/wb-go/wbf/rabbitmq"
)

// Publisher - обертка для отправки сообщений в RabbitMQ
type Publisher struct {
	client *rabbitmq.Publisher
}

// NewPublisher создает новый Publisher
func NewPublisher(client *rabbitmq.Publisher) *Publisher {
	return &Publisher{
		client: client,
	}
}

// Publish отправляет сообщение в очередь
func (p *Publisher) Publish(queueName string, data []byte) error {
	fmt.Println(queueName, "сообщение", data)
	return p.client.Publish(data, queueName, "application/json")
}

// PublishNotification отправляет уведомление в RabbitMQ
func (p *Publisher) PublishNotification(queueName string, notification map[string]interface{}) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	fmt.Println("Уведомление:", string(data))

	return p.Publish(queueName, data)
}
