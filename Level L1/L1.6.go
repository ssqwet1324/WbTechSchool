package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// printIntsCtx через контекст
func printIntsCtx(ctx context.Context, nums <-chan int, done chan<- bool) {
	for {
		select {
		case i := <-nums:
			time.Sleep(time.Millisecond * 100)
			fmt.Println("число:", i)
		case <-ctx.Done():
			fmt.Println(" конец контекста")
			done <- true
			return
		}
	}
}

// PrintIntsStop через канал уведомления
func PrintIntsStop(stop <-chan bool, nums <-chan int, done chan<- bool) {
	for {
		select {
		case <-stop:
			fmt.Println("конец канал уведомлений")
			done <- true
			return
		case i := <-nums:
			time.Sleep(time.Millisecond * 100)
			fmt.Println("число:", i)
		}
	}
}

// PrintInts по условию
func PrintInts(nums <-chan int, done chan<- bool) {
	for n := range nums {
		if n > 10 {
			fmt.Println("выход по условию")
			done <- true
			return
		}
		fmt.Println("число:", n)
	}
}

// PrintIntsExit - через runtime.Goexit
func PrintIntsExit(nums <-chan int, done chan<- bool) {
	defer func() {
		done <- true
	}()
	for n := range nums {
		time.Sleep(time.Millisecond * 100)
		if n == 10 {
			fmt.Println("выход через runtime.Goexit()")
			runtime.Goexit()
		}
		fmt.Println("число: ", n)
	}
}

func main() {
	// выход по контексту
	{
		ch := make(chan int, 50)
		done := make(chan bool)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		go printIntsCtx(ctx, ch, done)

		go func() {
			for i := 0; i < 1000; i++ {
				select {
				case <-ctx.Done():
					close(ch)
					return
				default:
					ch <- i
					time.Sleep(time.Millisecond)
				}
			}
		}()

		<-done
	}

	//выход по каналу уведомлений
	{
		ch := make(chan int, 50)
		stop := make(chan bool)
		done := make(chan bool)

		go PrintIntsStop(stop, ch, done)

		go func() {
			for i := 0; i < 1000; i++ {
				select {
				case <-stop:
					close(ch)
					return
				default:
					ch <- i
					time.Sleep(time.Millisecond)
				}
			}
		}()

		go func() {
			time.Sleep(time.Second)
			stop <- true
		}()

		<-done
	}

	// выход по условию
	{
		nums := make(chan int, 100)
		done := make(chan bool)

		go PrintInts(nums, done)

		go func() {
			for i := 0; i < 1000; i++ {
				nums <- i
				time.Sleep(time.Millisecond * 100)
			}
			close(nums)
		}()

		<-done
	}
	//выход через runtime.Goexit()
	{
		nums := make(chan int, 100)
		done := make(chan bool)
		go PrintIntsExit(nums, done)

		go func() {
			for i := 0; i < 1000; i++ {
				nums <- i
				time.Sleep(time.Millisecond * 100)
			}
			close(nums)
		}()
		<-done
	}
	fmt.Println("конец")
}
