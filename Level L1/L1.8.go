package main

import (
	"fmt"
)

func main() {
	var num int64
	var i int
	var bit int
	fmt.Println("Введите число")
	fmt.Scan(&num)
	fmt.Println("Введите номер бита")
	fmt.Scan(&i)
	fmt.Println("Введите 0 или 1, чтобы обнулить или установить бит:")
	fmt.Scan(&bit)

	if bit == 0 {
		num = num &^ (1 << i)
	} else {
		num = num | (1 << i)
	}

	fmt.Println(num)
}
