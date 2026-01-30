package main

import (
	"fmt"
	"strings"
)

func IsUnique(str string) bool {
	lower := strings.ToLower(str)
	mapa := make(map[rune]struct{})
	for _, v := range lower {
		if _, ok := mapa[v]; ok {
			return false
		}

		mapa[v] = struct{}{}
	}

	return true
}

func main() {
	fmt.Println(IsUnique("abCdeAf"))
}
