package analoguecut

import (
	"L2_13/internal/cli"
	"bufio"
	"strings"
)

// Cut - утилита cut
func Cut(scanner *bufio.Scanner, flags cli.Flags) string {
	// проверяем что в файле остались строки
	if !scanner.Scan() {
		return ""
	}

	line := scanner.Text()

	// Если задан -s пропускаем строки без разделителя
	if flags.Separator && !strings.Contains(line, flags.Delimiter) {
		return ""
	}

	// делим строки по разделителю
	splitLines := strings.Split(line, flags.Delimiter)
	var result []string

	// берем нужные слова
	for _, field := range flags.Fields {
		if field > 0 && field-1 < len(splitLines) {
			result = append(result, splitLines[field-1])
		}
	}

	// проверяем что список не пустой
	if len(result) == 0 {
		return ""
	}

	return strings.Join(result, flags.Delimiter)
}
