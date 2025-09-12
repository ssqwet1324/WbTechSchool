package search

import (
	"L2_12/internal/cli"
	"bufio"
	"fmt"
	"strings"
)

// PatternString - поиск строк по шаблону
func PatternString(scanner *bufio.Scanner, options *cli.Options, searchedStr string) {
	lineNum := 0                                      // номер строки
	afterCount := 0                                   // количество строк после шаблона для(-A)
	beforeBuffer := make([]string, 0, options.Before) // строки до шаблона для(-B)
	matches := 0                                      // количество найденных совпадений
	printedMatches := make(map[string]bool)           // чтобы не дублировать совпадения

	for scanner.Scan() {
		text := scanner.Text()
		lineNum++

		// // проверка для флага -i
		checkText := text
		checkStr := searchedStr
		if options.Ignore {
			checkText = strings.ToLower(checkText)
			checkStr = strings.ToLower(checkStr)
		}

		// проверка для флага -F
		matched := false
		if options.Fixed {
			matched = strings.Contains(checkText, checkStr)
		} else {
			matched = checkText == checkStr
		}

		// проверка для флага -v
		if options.Invert {
			matched = !matched
		}

		// // проверка для флага -c
		if matched {
			matches++
		}
		if options.Count {
			continue
		}

		switch {
		case matched:
			// выводим буфер перед совпадением (-B)
			for _, b := range beforeBuffer {
				fmt.Println(b)
			}

			// выводим совпадение, если его еще не выводили
			if !printedMatches[text] {
				printedMatches[text] = true
				if options.LineNum {
					fmt.Printf("%d:%s\n", lineNum, text)
				} else {
					fmt.Println(text)
				}
			}

			afterCount = options.After
			beforeBuffer = beforeBuffer[:0]

			// активируем -A
		case afterCount > 0:
			fmt.Println(text)
			afterCount--

			// тут активируем -B
		default:
			if options.Before > 0 {
				if len(beforeBuffer) == options.Before {
					beforeBuffer = beforeBuffer[1:]
				}
				beforeBuffer = append(beforeBuffer, text)
			}
		}
	}

	// считаем количество строк для шаблона -n
	if options.Count {
		fmt.Println(matches)
	}
}
