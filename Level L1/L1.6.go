package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// printIntsCtx — выход по контексту
func printIntsCtx(ctx context.Context, nums <-chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("конец контекста")
			return
		case i, ok := <-nums:
			if !ok {
				fmt.Println("nums closed")
				return
			}
			time.Sleep(100 * time.Millisecond)
			fmt.Println("число:", i)
		}
	}
}

// PrintIntsStop — через канал уведомления
func PrintIntsStop(stop <-chan struct{}, nums <-chan int) {
	for {
		select {
		case <-stop:
			fmt.Println("конец канал уведомлений")
			return
		case i, ok := <-nums:
			if !ok {
				fmt.Println("nums closed")
				return
			}
			time.Sleep(100 * time.Millisecond)
			fmt.Println("число:", i)
		}
	}
}

// PrintInts — выход по условию
func PrintInts(nums <-chan int) {
	for n := range nums {
		if n > 10 {
			fmt.Println("выход по условию")
			return
		}
		fmt.Println("число:", n)
	}
}

// PrintIntsExit — выход через runtime.Goexit
func PrintIntsExit(nums <-chan int) {
	for n := range nums {
		time.Sleep(100 * time.Millisecond)
		if n == 10 {
			fmt.Println("выход через runtime.Goexit()")
			runtime.Goexit()
		}
		fmt.Println("число:", n)
	}
}

func main() {
	// выход по контексту
	{
		ch := make(chan int, 50)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			printIntsCtx(ctx, ch)
		}()

		go func() {
			defer close(ch)
			for i := 0; i < 1000; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
					time.Sleep(time.Millisecond)
				}
			}
		}()

		wg.Wait()
	}

	// выход по каналу уведомлений
	{
		ch := make(chan int, 50)
		stop := make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			PrintIntsStop(stop, ch)
		}()

		go func() {
			defer close(ch)
			for i := 0; i < 1000; i++ {
				select {
				case <-stop:
					return
				case ch <- i:
					time.Sleep(time.Millisecond)
				}
			}
		}()

		go func() {
			time.Sleep(time.Second)
			close(stop)
		}()

		wg.Wait()
	}

	// выход по условию
	{
		nums := make(chan int, 100)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			PrintInts(nums)
		}()

		go func() {
			defer close(nums)
			for i := 0; i < 1000; i++ {
				nums <- i
				time.Sleep(100 * time.Millisecond)
			}
		}()

		wg.Wait()
	}

	// выход через runtime.Goexit()
	{
		nums := make(chan int, 100)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			PrintIntsExit(nums)
		}()

		go func() {
			defer close(nums)
			for i := 0; i < 1000; i++ {
				nums <- i
				time.Sleep(100 * time.Millisecond)
			}
		}()

		wg.Wait()
	}

	fmt.Println("конец")
}
