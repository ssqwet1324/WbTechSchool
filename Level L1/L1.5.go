package main

import (
	"fmt"
	"time"
)

// Send - отправляем данные в канал
func Send() <-chan int {
	ch := make(chan int)
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()
		defer close(ch)

		for i := 0; i < 1000; i++ {
			<-ticker.C
			ch <- i
		}
	}()

	return ch
}

func main() {
	ch := Send()

	var seconds int
	fmt.Println("Введите количество секунд")
	if _, err := fmt.Scan(&seconds); err != nil {
		panic("Введен некорректный символ")
	}

	timeout := time.After(time.Duration(seconds) * time.Second)

	//тут берем данные с канала и выводим пока таймер не вышел
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println(val)
		case <-timeout:
			fmt.Println("время вышло")
			return
		}
	}
}
