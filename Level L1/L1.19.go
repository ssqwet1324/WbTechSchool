package main

import "fmt"

func ReverseStr(str string) {
	runes := []rune(str)
	for i := len(runes) - 1; i >= 0; i-- {
		fmt.Printf("%c", runes[i])
	}
	fmt.Print("\n")
}

func main() {
	var symbol string
	if _, err := fmt.Scan(&symbol); err != nil {
		panic("введено не корректное значение")
	}

	ReverseStr(symbol)
}
