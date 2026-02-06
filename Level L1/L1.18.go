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

func (c *Counter) Increment(count int) {
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			c.mu.Lock()
			c.value++
			v := c.value
			c.mu.Unlock()

			time.Sleep(100 * time.Millisecond) // для имитации работы
			fmt.Println("Increment", v)
		}()
	}

	wg.Wait()
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	counter := &Counter{}
	countGoroutines := 100

	counter.Increment(countGoroutines)
	fmt.Println("Counter Value:", counter.Value())
}
