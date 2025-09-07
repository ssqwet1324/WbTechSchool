package main

import (
	"fmt"
	"slices"
)

// groupAnagrams - нахождение анаграмм
func groupAnagrams(strings []string) [][]string {
	//создаем мапу для анаграмм
	anagramsMap := make(map[string][]string, len(strings))
	//проходимся по строкам сортируем и создаем ключ со значениями
	for _, str := range strings {
		chars := []rune(str)
		slices.Sort(chars)
		key := string(chars)
		anagramsMap[key] = append(anagramsMap[key], str)
	}

	//выводим значения
	var anagrams [][]string
	for _, anagram := range anagramsMap {
		anagrams = append(anagrams, anagram)
	}

	return anagrams
}

func main() {
	anagrams := groupAnagrams([]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"})

	for i := 0; i < len(anagrams); i++ {
		if len(anagrams[i]) != 1 {
			fmt.Printf("%v: %v\n", anagrams[i][0], anagrams[i])
		}
	}
}
