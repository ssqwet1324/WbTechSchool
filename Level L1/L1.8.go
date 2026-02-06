package main

import (
	"fmt"
)

func main() {
	var num int64
	var i int
	var bit int

	fmt.Print("Введите число: ")
	if _, err := fmt.Scan(&num); err != nil {
		panic("Введено некорректное значение")
	}

	fmt.Print("Введите номер бита (0..63): ")
	if _, err := fmt.Scan(&i); err != nil {
		panic("Введено некорректное значение")
	}
	if i < 0 || i > 63 {
		panic("Номер бита должен быть в диапазоне 0..63")
	}

	fmt.Print("Введите 0 или 1: ")
	if _, err := fmt.Scan(&bit); err != nil {
		panic("Введено некорректное значение")
	}
	if bit != 0 && bit != 1 {
		panic("Нужно ввести 0 или 1")
	}

	if bit == 0 {
		num = num &^ (1 << i)
	} else {
		num = num | (1 << i)
	}

	fmt.Println(num)
}
