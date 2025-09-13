package cli

import (
	"flag"
	"log"
	"strconv"
	"strings"
)

// Flags - структура хранящая флаги
type Flags struct {
	Fields    []int
	Delimiter string
	Separator bool
}

// ParseFlags - парсим флаги
func ParseFlags() (Flags, string) {
	var delimiter, fields string
	var separated bool
	var parsedFields []int

	flag.StringVar(&fields, "f", "", "Номера полей которые нужно вывести")
	flag.StringVar(&delimiter, "d", "\t", "Разделитель")
	flag.BoolVar(&separated, "s", false, "Показывать только строки содержащие разделитель")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("нужно передать имя файла")
	}
	fileName := args[0]

	parsedFields = ParseFields(fields)

	flags := Flags{
		Delimiter: delimiter,
		Fields:    parsedFields,
		Separator: separated,
	}

	return flags, fileName
}

// ParseFields - парсим , и  -
func ParseFields(s string) []int {
	var result []int

	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			// Диапазон
			bounds := strings.SplitN(part, "-", 2)
			start, err1 := strconv.Atoi(strings.TrimSpace(bounds[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err1 != nil || err2 != nil || start > end {
				continue // некорректный диапазон игнорируем
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// Одиночное число
			if n, err := strconv.Atoi(part); err == nil {
				result = append(result, n)
			}
		}
	}

	return result
}
