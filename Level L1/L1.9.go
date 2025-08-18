package main

import "fmt"

// берем числа из списка с помощью горутины
func GiveNums(nums ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, num := range nums {
			c <- num
		}
		close(c)
	}()
	return c
}

// получаем числа с канала и умножаем
func SqrtNums(num <-chan int) <-chan int {
	c := make(chan int)

	go func() {
		for n := range num {
			c <- n * 2
		}
		close(c)
	}()
	return c
}

// вызываем
func main() {
	for n := range SqrtNums(GiveNums(2, 4, 6, 8, 10, 12, 14, 16, 18, 20)) {
		fmt.Println(n)
	}
}
