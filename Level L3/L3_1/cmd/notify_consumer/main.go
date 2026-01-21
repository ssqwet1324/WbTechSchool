package main

import (
	"context"
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
)

const (
	rabbitLocalUrl = "amqp://guest:guest@localhost:5672/"
	retries        = 5
	pause          = 2 * time.Second
)

// консюмер
func main() {
	conn, err := rabbitmq.Connect(rabbitLocalUrl, retries, pause)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	qm := rabbitmq.NewQueueManager(ch)
	queue, err := qm.DeclareQueue("notify_queue")
	if err != nil {
		log.Fatalf("declare queue: %v", err)
	}
	_ = queue

	consumer := rabbitmq.NewConsumer(ch, rabbitmq.NewConsumerConfig("notify_queue"))
	msgChan := make(chan []byte)

	go func() {
		if err := consumer.Consume(msgChan); err != nil {
			log.Printf("consume error: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("rabbitmq started: waiting for messages from notify_queue...")
	for {
		select {
		case <-ctx.Done():
			return
		case body := <-msgChan:
			log.Printf("received: %s", string(body))
		}
	}
}
