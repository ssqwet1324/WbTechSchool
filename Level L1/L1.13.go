package main

import "fmt"

func main() {
	var num1, num2 int
	_, err := fmt.Scan(&num1, &num2)
	if err != nil {
		panic("Введено не число")
	}
	num2 = num1 + num2
	num1 = num2 - num1
	num2 = num2 - num1

	fmt.Println(num1, num2)
}
