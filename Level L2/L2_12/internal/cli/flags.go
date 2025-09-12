package cli

import (
	"flag"
	"log"
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

func ParseOptions() (Options, string, string) {
	// Флаги
	flagA := flag.Int("A", 0, "вывести N строк после найденной строки")
	flagB := flag.Int("B", 0, "вывести N строк до найденной строки")
	flagC := flag.Int("C", 0, "вывести N строк контекста вокруг найденной строки")
	flagCCount := flag.Bool("c", false, "выводить только количество совпадающих строк")
	flagI := flag.Bool("i", false, "игнорировать регистр")
	flagV := flag.Bool("v", false, "инвертировать фильтр")
	flagF := flag.Bool("F", false, "воспринимать шаблон как фиксированную строку")
	flagN := flag.Bool("n", false, "выводить номер строки перед каждой найденной строкой")

	flag.Parse()

	// Проверка аргументов
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("нужно передать шаблон и имя файла")
	}
	pattern := args[0]
	filename := args[1]

	after, before := *flagA, *flagB
	if *flagC > 0 {
		// Превращаем -С в -А и -В
		after = *flagC
		before = *flagC
	}

	return Options{
		After:   after,
		Before:  before,
		Count:   *flagCCount,
		Ignore:  *flagI,
		Invert:  *flagV,
		Fixed:   *flagF,
		LineNum: *flagN,
	}, pattern, filename
}
