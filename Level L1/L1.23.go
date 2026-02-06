package main

import "fmt"

func main() {
	var n int
	fmt.Print("Введите размер слайса: ")
	if _, err := fmt.Scan(&n); err != nil || n < 0 {
		panic("некорректный размер")
	}

	sp := make([]int, n)
	fmt.Println("Введите элементы:")
	for i := 0; i < n; i++ {
		if _, err := fmt.Scan(&sp[i]); err != nil {
			panic("некорректный элемент")
		}
	}

	var elem int
	fmt.Print("Введите индекс, который удалить (0..n-1): ")
	if _, err := fmt.Scan(&elem); err != nil {
		panic("некорректный индекс")
	}
	if elem < 0 || elem >= len(sp) {
		panic("индекс вне диапазона")
	}

	// 1 способ: через copy (удаляет элемент по индексу)
	spCopy := append([]int(nil), sp...)
	copy(spCopy[elem:], spCopy[elem+1:])
	spCopy = spCopy[:len(spCopy)-1]
	fmt.Println("copy:", spCopy)

	// 2 способ: через append
	spAppend := append([]int(nil), sp...)
	spAppend = append(spAppend[:elem], spAppend[elem+1:]...)
	fmt.Println("append:", spAppend)
}
