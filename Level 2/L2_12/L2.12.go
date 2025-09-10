package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Options - структура для хранения данных флагов
type Options struct {
	After   int  // -A
	Before  int  // -B
	Count   bool // -c
	Ignore  bool // -i
	Invert  bool // -v
	Fixed   bool // -F
	LineNum bool // -n
}

// SearchString - функция для нахождения строки
func SearchString(scanner *bufio.Scanner, options *Options, searchedStr string) {
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

func main() {
	// указываем флаги
	flagA := flag.Int("A", 0, "вывести N строк после неё")
	flagB := flag.Int("B", 0, "вывести N строк до каждой найденной строки")
	flagC := flag.Int("C", 0, "вывести N строк контекста вокруг найденной строки")
	flagCCount := flag.Bool("c", false, "выводить только количество совпадающих строк")
	flagI := flag.Bool("i", false, "игнорировать регистр")
	flagV := flag.Bool("v", false, "инвертировать фильтр: выводить строки, не содержащие шаблон")
	flagF := flag.Bool("F", false, "воспринимать шаблон как фиксированную строку")
	flagN := flag.Bool("n", false, "выводить номер строки перед каждой найденной строкой")
	flag.Parse()

	// берем строку из терминала
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("нужно передать шаблон")
	}

	// берем имя файла из терминала
	var filename string
	if len(args) > 1 {
		filename = args[1]
	}

	// инициализируем структуру
	flagsStruct := Options{
		After:   *flagA,
		Before:  *flagB,
		Count:   *flagCCount,
		Ignore:  *flagI,
		Invert:  *flagV,
		Fixed:   *flagF,
		LineNum: *flagN,
	}

	// тут проверяем, если флаг -C активен то вызываем -A -B
	if *flagC > 0 {
		flagsStruct.After = *flagC
		flagsStruct.Before = *flagC
	}

	// проверка на пустую строку
	searchedString := args[0]
	if len(searchedString) == 0 {
		log.Fatal("передаваемая строка пустая")
	}

	//читаем или из файла, или из Stdin
	var scanner *bufio.Scanner
	if filename == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("ошибка открытия файла", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	// вызываем функцию поиска
	SearchString(scanner, &flagsStruct, searchedString)

	// тут обрабатываем ошибку чтения файла
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка чтения:", err)
	}
}
