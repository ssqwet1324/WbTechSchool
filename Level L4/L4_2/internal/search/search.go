package search

import (
	"fmt"
	"strings"

	"L4.2/internal/entity"
)

// SearchResult содержит результаты поиска
type SearchResult struct {
	Lines []string
	Count int
}

// SearchLines выполняет поиск с учётом всех флагов grep
func SearchLines(lines []string, options entity.Options, pattern string, offset int) (SearchResult, error) {
	if len(lines) == 0 {
		return SearchResult{Lines: nil, Count: 0}, nil
	}

	// Флаг -c: только количество
	if options.Count {
		count := countMatches(lines, pattern, options)
		return SearchResult{Lines: nil, Count: count}, nil
	}

	// Обычный поиск с контекстом
	result := searchWithContext(lines, pattern, options, offset)
	return SearchResult{Lines: result, Count: 0}, nil
}

// countMatches считает количество совпадений
func countMatches(lines []string, pattern string, options entity.Options) int {
	count := 0
	for _, line := range lines {
		if matches(line, pattern, options) {
			count++
		}
	}
	return count
}

// searchWithContext выполняет поиск с учётом флагов
func searchWithContext(lines []string, pattern string, options entity.Options, offset int) []string {
	if pattern == "" && !options.Invert {
		// Если паттерн пустой и не инвертирован, возвращаем все строки
		result := make([]string, 0, len(lines))
		for i, line := range lines {
			lineNum := offset + i + 1
			if options.LineNum {
				result = append(result, fmt.Sprintf("%d:%s", lineNum, line))
			} else {
				result = append(result, line)
			}
		}
		return result
	}

	result := make([]string, 0)
	printed := make(map[int]bool) // чтобы не дублировать строки

	// Буфер для хранения строк до совпадения (-B)
	type lineWithNum struct {
		line string
		num  int
	}
	beforeBuffer := make([]lineWithNum, 0, options.Before)
	afterCount := 0 // счётчик строк после совпадения (-A)

	for i, line := range lines {
		lineNum := offset + i + 1
		matched := matches(line, pattern, options)

		if matched {
			// Выводим буфер перед совпадением (-B)
			for _, bufItem := range beforeBuffer {
				if !printed[bufItem.num] {
					result = append(result, formatLine(bufItem.line, bufItem.num-1, options))
					printed[bufItem.num] = true
				}
			}
			beforeBuffer = beforeBuffer[:0]

			// Выводим само совпадение
			if !printed[lineNum] {
				result = append(result, formatLine(line, lineNum-1, options))
				printed[lineNum] = true
			}

			// Активируем счётчик для вывода строк после (-A)
			afterCount = options.After
		} else {
			// Если активен счётчик после совпадения, выводим строку
			if afterCount > 0 {
				if !printed[lineNum] {
					result = append(result, formatLine(line, lineNum-1, options))
					printed[lineNum] = true
				}
				afterCount--
			}

			// Обновляем буфер перед совпадением (-B)
			if options.Before > 0 {
				if len(beforeBuffer) >= options.Before {
					beforeBuffer = beforeBuffer[1:]
				}
				beforeBuffer = append(beforeBuffer, lineWithNum{line: line, num: lineNum})
			}
		}
	}

	return result
}

// matches проверяет, соответствует ли строка паттерну с учётом флагов
func matches(line, pattern string, options entity.Options) bool {
	if pattern == "" {
		return options.Invert // если паттерн пустой, -v вернёт все строки
	}

	checkLine := line
	checkPattern := pattern

	// Флаг -i: игнорировать регистр
	if options.Ignore {
		checkLine = strings.ToLower(line)
		checkPattern = strings.ToLower(pattern)
	}

	// Флаг -F: фиксированная строка (обычный поиск подстроки)
	var matched bool
	if options.Fixed {
		matched = strings.Contains(checkLine, checkPattern)
	} else {
		// По умолчанию тоже поиск подстроки (как в оригинальном grep)
		matched = strings.Contains(checkLine, checkPattern)
	}

	// Флаг -v: инвертировать результат
	if options.Invert {
		matched = !matched
	}

	return matched
}

// formatLine форматирует строку для вывода с учётом флага -n
func formatLine(line string, lineNum int, options entity.Options) string {
	if options.LineNum {
		return fmt.Sprintf("%d:%s", lineNum+1, line)
	}
	return line
}

// FilterLines - обратная совместимость (старая функция)
func FilterLines(lines []string, pattern string) []string {
	options := entity.Options{}
	result, _ := SearchLines(lines, options, pattern, 0)
	return result.Lines
}
