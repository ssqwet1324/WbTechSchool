package main

import (
	"fmt"
)

// UniqueStrings слайс уникальных строк
func UniqueStrings(input []string) []string {
	set := make(map[string]struct{}, len(input))
	for _, v := range input {
		set[v] = struct{}{}
	}

	res := make([]string, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	return res
}

func main() {
	sp1 := []string{"cat", "cat", "cat", "dog", "cat", "tree"}
	fmt.Println(UniqueStrings(sp1))

	sp2 := []string{"a", "b", "a", "c", "b", "d"}
	fmt.Println(UniqueStrings(sp2))

	var sp3 []string
	fmt.Println(UniqueStrings(sp3))
}
