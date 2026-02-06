package main

import "fmt"

func main() {
	nums := []int{2, 4, 6, 8, 10}
	done := make(chan struct{})

	for _, n := range nums {
		go func(x int) {
			fmt.Println(x * x)
			done <- struct{}{}
		}(n)
	}

	// ждём, пока все горутины закончат
	for range nums {
		<-done
	}

	close(done)
}
