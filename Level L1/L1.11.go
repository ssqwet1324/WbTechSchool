package main

import "fmt"

func main() {
	sp1 := []int{1, 2, 3}
	sp2 := []int{2, 3, 4}
	var result []int

	mn := make(map[int]struct{})

	//создаем ключ
	for _, v := range sp1 {
		mn[v] = struct{}{}
	}

	//проверяем есть ли элемент в мапе
	for _, v := range sp2 {
		//если есть то пересекаются
		if _, ok := mn[v]; ok {
			result = append(result, v)
		}
	}

	fmt.Println(result)
}
