package main

import "fmt"

func main() {
	sp1 := []string{"cat", "cat", "cat", "dog", "cat", "tree"}
	mapa := make(map[string]struct{})
	var result []string

	//создаем ключ(т.к они не повторяются)
	for _, v := range sp1 {
		//как заглушка
		mapa[v] = struct{}{}
	}

	for k, _ := range mapa {
		result = append(result, k)
	}

	fmt.Println(result)
}
