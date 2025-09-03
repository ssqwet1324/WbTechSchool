package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// UnpackingString - распаковка строк
func UnpackingString(str string) (string, error) {
	// пустая строка
	if str == "" {
		return "", nil
	}

	// проверка что только цифры
	allDigits := true
	for _, r := range str {
		if !unicode.IsDigit(r) {
			allDigits = false
			break
		}
	}
	if allDigits {
		return "", errors.New("некорректная строка: только цифры")
	}

	runes := []rune(str)
	var result []rune
	var pred rune //предыдущий символ

	for i := 0; i < len(runes); i++ {
		r := runes[i] //текущий символ

		if unicode.IsDigit(r) {
			// тут проверка на строку такого типа: 4a
			if pred == 0 {
				return "", errors.New("некорректная строка: цифра без символа")
			}
			num, _ := strconv.Atoi(string(r))
			for j := 0; j < num-1; j++ {
				result = append(result, pred)
			}
		} else {
			result = append(result, r)
			pred = r
		}
	}

	return string(result), nil
}

func main() {
	tests := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		"4a",
	}

	for _, t := range tests {
		res, err := UnpackingString(t)
		fmt.Println("ввод:", t, "вывод: ", res, err)
	}
}
