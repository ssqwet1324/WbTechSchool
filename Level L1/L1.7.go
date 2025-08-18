package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const workerCounts = 3

func Writer(ctx context.Context, id int, num <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-num:
			if !ok {
				return
			}
			fmt.Println("Worker", id, "write", n)
			time.Sleep(time.Second)
			results <- n
		}
	}
}

func main() {
	mapa := make(map[string]int)
	mu := sync.RWMutex{}
	num := make(chan int, 50)
	results := make(chan int, 50)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	//wg для того чтобы удобно закрыть канал result и writer
	var wg sync.WaitGroup
	wg.Add(workerCounts)
	for i := 0; i < workerCounts; i++ {
		go Writer(ctx, i, num, results, &wg)
	}

	// пишем безопасно в мапу
	go func() {
		for n := range results {
			mu.Lock()
			mapa[strconv.Itoa(n)] = n
			mu.Unlock()
		}
	}()

	// отправляем числа в канал
	go func() {
		n := 1
		for {
			select {
			case <-ctx.Done():
				close(num)
				return
			case num <- n:
				n++
			}
		}
	}()

	//ожидаем завершения горутин
	wg.Wait()
	//закрываем канал
	close(results)
	//тут безопасно выводим мапу
	mu.RLock()
	defer mu.RUnlock()
	for k, v := range mapa {
		fmt.Println("key", k, "value", v)
	}
}
