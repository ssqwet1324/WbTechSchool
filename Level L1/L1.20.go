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

func ReverseWordsRune(str string) {
	runes := []rune(str)

	//разворачиваем всю строку
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	//разворачиваем каждое слово отдельно
	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || runes[i] == ' ' {
			for l, r := start, i-1; l < r; l, r = l+1, r-1 {
				runes[l], runes[r] = runes[r], runes[l]
			}
			start = i + 1
		}
	}

	fmt.Println(string(runes))
}

func main() {
	ReverseWords("snow dog sun")
	ReverseWordsRune("snow dog sun")
}
