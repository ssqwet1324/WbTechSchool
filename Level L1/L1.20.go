package main

import (
	"fmt"
	"strings"
)

func ReverseWords(str string) {
	words := strings.Fields(str)
	for i := len(words) - 1; i >= 0; i-- {
		fmt.Print(words[i] + " ")
	}
}

func main() {
	ReverseWords("snow dog sun")
}
