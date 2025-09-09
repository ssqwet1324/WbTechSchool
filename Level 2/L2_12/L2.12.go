package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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
	lineNum := 0
	afterCount := 0
	beforeBuffer := make([]string, 0, options.Before)

	// проходимся по строкам в файле
	for scanner.Scan() {
		text := scanner.Text()
		lineNum++

		// проверка совпала строка или нет
		matched := text == searchedStr

		switch {
		//если совпала строка
		case matched:
			// печатаем строки из среза для флага -B
			if options.Before > 0 {
				for _, b := range beforeBuffer {
					fmt.Println(b)
				}
			}

			fmt.Println(text)

			// включаем -A
			afterCount = options.After

			// сбрасываем буфер
			beforeBuffer = beforeBuffer[:0]

		// тут допилить количетсво строк совпадающиъ

		// выводим n строк после совпадения для -A
		case afterCount > 0:
			fmt.Println(text)
			afterCount--

		// накапливаем строки для -B
		default:
			if options.Before > 0 {
				if len(beforeBuffer) == options.Before {
					beforeBuffer = beforeBuffer[1:] // удаляем старую строку
				}
				beforeBuffer = append(beforeBuffer, text)
			}
		}
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
			log.Fatal(err)
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
