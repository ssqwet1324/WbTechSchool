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
	ReverseStr("главрыба")
	ReverseStr("главрыба👉")
	ReverseStr("👉👈🧠🍒")
}
