package cli

import (
	"flag"

	"L4.2/internal/entity"
)

var (
	flagA      = flag.Int("A", 0, "вывести N строк после найденной строки")
	flagB      = flag.Int("B", 0, "вывести N строк до найденной строки")
	flagC      = flag.Int("C", 0, "вывести N строк контекста вокруг найденной строки")
	flagCCount = flag.Bool("c", false, "выводить только количество совпадающих строк")
	flagI      = flag.Bool("i", false, "игнорировать регистр")
	flagV      = flag.Bool("v", false, "инвертировать фильтр")
	flagF      = flag.Bool("F", false, "воспринимать шаблон как фиксированную строку")
	flagN      = flag.Bool("n", false, "выводить номер строки перед каждой найденной строкой")
)

// ParseOptions - парсим флаги для grep
func ParseOptions() entity.Options {
	flag.Parse()

	after, before := *flagA, *flagB
	if *flagC > 0 {
		after = *flagC
		before = *flagC
	}

	return entity.Options{
		After:   after,
		Before:  before,
		Count:   *flagCCount,
		Ignore:  *flagI,
		Invert:  *flagV,
		Fixed:   *flagF,
		LineNum: *flagN,
	}
}
