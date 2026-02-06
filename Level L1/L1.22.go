package main

import (
	"fmt"
	"math/big"
)

func main() {
	var num1, num2 big.Int

	if _, err := fmt.Scan(&num1, &num2); err != nil {
		panic("введено не число")
	}

	sum := new(big.Int).Add(&num1, &num2)
	diff := new(big.Int).Sub(&num1, &num2)
	prod := new(big.Int).Mul(&num1, &num2)

	fmt.Printf("Сумма: %v\nРазность: %v\nПроизведение: %v\n", sum, diff, prod)

	if num2.Sign() == 0 {
		fmt.Println("Деление: нельзя делить на 0")
		return
	}
	quot := new(big.Int).Div(&num1, &num2)
	fmt.Printf("Деление: %v\n", quot)
}
