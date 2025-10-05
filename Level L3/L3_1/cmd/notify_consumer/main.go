package main

import (
	"context"
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
)

// Пример простого консюмера, который читает сообщения из очереди notify_queue и печатает их в stdout.
func main() {
	conn, err := rabbitmq.Connect("amqp://guest:guest@localhost:5672/", 5, 2*time.Second)
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
	_ = queue // очередь объявлена

	consumer := rabbitmq.NewConsumer(ch, rabbitmq.NewConsumerConfig("notify_queue"))
	msgChan := make(chan []byte)

	go func() {
		if err := consumer.Consume(msgChan); err != nil {
			log.Printf("consume error: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("consumer started: waiting for messages from notify_queue...")
	for {
		select {
		case <-ctx.Done():
			return
		case body := <-msgChan:
			log.Printf("received: %s", string(body))
		}
	}
}
