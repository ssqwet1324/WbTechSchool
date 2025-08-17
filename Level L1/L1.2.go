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

// получаем числа с канала и возводим в квадрат
func SqrtNums(num <-chan int) <-chan int {
	c := make(chan int)

	go func() {
		for n := range num {
			c <- n * n
		}
		close(c)
	}()
	return c
}

// вызываем
func main() {
	for n := range SqrtNums(GiveNums(2, 4, 6, 8, 10)) {
		fmt.Println(n)
	}
}
