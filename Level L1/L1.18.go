package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment(count int, stop chan bool) {
	for i := 0; i < count; i++ {
		go func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.value++
			time.Sleep(time.Millisecond * 100)
			fmt.Println("Increment", c.value)
			stop <- true
		}()
	}

}

func (c *Counter) Value() int {
	return c.value
}

func main() {
	counter := &Counter{mu: sync.Mutex{}, value: 0}
	countGoroutines := 100
	stop := make(chan bool, countGoroutines)

	counter.Increment(countGoroutines, stop)

	for i := 0; i < countGoroutines; i++ {
		<-stop
	}

	fmt.Println("Counter Value:", counter.Value())
}
