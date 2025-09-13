package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	chDone := make(chan interface{})

	// запускаем отдельную горутину
	go func() {
		// проходимся по каналам
		for _, channel := range channels {
			go func() {
				// ждем пока в канал не отправится значение или закроется
				<-channel
				select {
				// отправляем значение
				case chDone <- struct{}{}:
				default:
				}
			}()
		}
	}()

	return chDone
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
