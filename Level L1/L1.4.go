package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const workersCount = 5

func workers(ctx context.Context, id int, works <-chan int, results chan<- int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker ", id, " stopped")
			return
		case n := <-works:
			fmt.Println("Worker", id, "start", n)
			time.Sleep(time.Millisecond * time.Duration(n))
			fmt.Println("Worker", id, "end", n)
			results <- n
		}
	}
}

func main() {
	works := make(chan int, 50)
	results := make(chan int, 50)

	//создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем 5 воркеров читающих результат
	for i := 1; i <= workersCount; i++ {
		go workers(ctx, i, works, results)
	}

	//создаем канал для отлавливания сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	//выводим что прочитал воркер
	go func() {
		for n := range results {
			fmt.Println(n)
		}
	}()

	// постоянно пишем
	go func() {
		workID := 1
		for {
			select {
			case <-ctx.Done():
				close(works)
				return
			default:
				works <- workID
				workID++
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()
	//тут слушаем сигнал
	<-signalChan
	//даем сигнал воркерам что пора завершаться
	cancel()
	time.Sleep(time.Second * 1)
	close(results)
	fmt.Println("Done")
}
