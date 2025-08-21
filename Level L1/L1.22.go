package main

import (
	"fmt"
	"math/big"
)

func main() {
	a, _ := new(big.Int).SetString("123456789123456789123456789", 10)
	b, _ := new(big.Int).SetString("2400000000000000000000", 10)

	sum := new(big.Int).Add(a, b)
	diff := new(big.Int).Sub(a, b)
	prod := new(big.Int).Mul(a, b)
	quot := new(big.Int).Div(a, b)

	fmt.Printf("Сумма: %v\nРазность: %v\nПроизведение: %v\nДеление: %v\n", sum, diff, prod, quot)
}
