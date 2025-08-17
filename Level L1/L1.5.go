package main

import (
	"fmt"
	"time"
)

// отправляем данные в канал
func Send() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 1000; i++ {
			ch <- i
			time.Sleep(time.Millisecond)
		}
		close(ch)
	}()
	return ch
}

func main() {
	ch := Send()
	timeout := time.After(time.Second)

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
