//package main
//
//import (
//	"bufio"
//	"flag"
//	"fmt"
//	"io"
//	"log"
//	"os"
//	"sort"
//	"strconv"
//	"strings"
//)
//
//
//const (
//	chunkSize = 1000
//	filename = "test.txt"
//)
//
//// GetColumn - получаем колонку из строки
//func GetColumn(line string, n int) string {
//	cols := strings.Split(line, "\t")
//	if n-1 < len(cols) {
//		return cols[n-1]
//	}
//	return ""
//}
//
//// SortLines - сортировка среза строк в памяти
//func SortLines(lines []string, column int, numeric, reverse bool) {
//	sort.Slice(lines, func(i, j int) bool {
//		col1 := GetColumn(lines[i], column)
//		col2 := GetColumn(lines[j], column)
//
//		if numeric {
//			ai, errA := strconv.Atoi(col1)
//			bi, errB := strconv.Atoi(col2)
//			if errA == nil && errB == nil {
//				if reverse {
//					return ai > bi
//				}
//				return ai < bi
//			}
//			// Если не удалось преобразовать в числа, сортируем как строки
//		}
//
//		if reverse {
//			return col1 > col2
//		}
//		return col1 < col2
//	})
//}
//
//// ReadChunk - чтение куска строк из сканера
//func ReadChunk(scanner *bufio.Scanner, chunkSize int) ([]string, bool) {
//	lines := make([]string, 0, chunkSize)
//	for i := 0; i < chunkSize && scanner.Scan(); i++ {
//		lines = append(lines, scanner.Text())
//	}
//	return lines, len(lines) > 0
//}
//
//// ExternalSort - сортировка кусками
//func ExternalSort(input io.Reader, output io.Writer, column int, numeric, reverse, unique bool) error {
//	scanner := bufio.NewScanner(input)
//	var chunkFiles []string
//
//	// Создаём временные отсортированные файлы
//	for chunkNum := 0; ; chunkNum++ {
//		lines, ok := ReadChunk(scanner, chunkSize)
//		if !ok {
//			break
//		}
//
//		SortLines(lines, column, numeric, reverse)
//
//		tmpFile, err := os.CreateTemp("", fmt.Sprintf("chunk_%d_*.txt", chunkNum))
//		if err != nil {
//			return err
//		}
//
//		writer := bufio.NewWriter(tmpFile)
//		for _, line := range lines {
//			fmt.Fprintln(writer, line)
//		}
//		writer.Flush()
//		tmpFile.Close()
//		chunkFiles = append(chunkFiles, tmpFile.Name())
//	}
//
//	if len(chunkFiles) == 0 {
//		return nil
//	}
//
//	// Слияние временных файлов
//	type fileLine struct {
//		line    string
//		scanner *bufio.Scanner
//		file    *os.File
//	}
//
//	openFiles := make([]fileLine, 0, len(chunkFiles))
//	for _, fname := range chunkFiles {
//		f, err := os.Open(fname)
//		if err != nil {
//			return err
//		}
//		sc := bufio.NewScanner(f)
//		if sc.Scan() {
//			openFiles = append(openFiles, fileLine{
//				line:    sc.Text(),
//				scanner: sc,
//				file:    f,
//			})
//		} else {
//			f.Close()
//		}
//	}
//
//	var lastPrinted string
//	writer := bufio.NewWriter(output)
//	defer writer.Flush()
//
//	for len(openFiles) > 0 {
//		// Находим следующую строку для вывода (мин или макс в зависимости от reverse)
//		selectedIdx := 0
//		for i := 1; i < len(openFiles); i++ {
//			col1 := GetColumn(openFiles[i].line, column)
//			col2 := GetColumn(openFiles[selectedIdx].line, column)
//
//			var shouldSwap bool
//
//			if numeric {
//				num1, err1 := strconv.Atoi(col1)
//				num2, err2 := strconv.Atoi(col2)
//				if err1 == nil && err2 == nil {
//					// Оба числовые
//					if reverse {
//						shouldSwap = num1 > num2 // для обратной сортировки ищем максимум
//					} else {
//						shouldSwap = num1 < num2 // для прямой сортировки ищем минимум
//					}
//				} else {
//					// Хотя бы одно не числовое - сортируем как строки
//					if reverse {
//						shouldSwap = col1 > col2
//					} else {
//						shouldSwap = col1 < col2
//					}
//				}
//			} else {
//				// Строковая сортировка
//				if reverse {
//					shouldSwap = col1 > col2
//				} else {
//					shouldSwap = col1 < col2
//				}
//			}
//
//			if shouldSwap {
//				selectedIdx = i
//			}
//		}
//
//		// Проверка на уникальность и вывод
//		currentLine := openFiles[selectedIdx].line
//		if !unique || currentLine != lastPrinted {
//			fmt.Fprintln(writer, currentLine)
//			lastPrinted = currentLine
//		}
//
//		// Читаем следующую строку из выбранного файла
//		if openFiles[selectedIdx].scanner.Scan() {
//			openFiles[selectedIdx].line = openFiles[selectedIdx].scanner.Text()
//		} else {
//			// Файл закончился - закрываем и удаляем из списка
//			openFiles[selectedIdx].file.Close()
//			openFiles = append(openFiles[:selectedIdx], openFiles[selectedIdx+1:]...)
//		}
//	}
//
//	// Удаляем временные файлы
//	for _, fname := range chunkFiles {
//		os.Remove(fname)
//	}
//
//	return nil
//}
//
//func main() {
//	// Определяем флаги
//	k := flag.Int("k", 1, "sort by column number")
//	n := flag.Bool("n", false, "numeric sort")
//	r := flag.Bool("r", false, "reverse sort")
//	u := flag.Bool("u", false, "unique lines only")
//
//	// Также поддерживаем комбинированные флаги типа -nr
//	flag.BoolVar(n, "nr", false, "numeric reverse sort (combined)")
//	flag.BoolVar(r, "rn", false, "reverse numeric sort (combined)")
//	flag.BoolVar(n, "nru", false, "numeric reverse unique sort (combined)")
//
//	flag.Parse()
//
//
//	file, err := os.Open(filename)
//	if err != nil {
//		log.Fatalf("Error opening file: %v", err)
//	}
//	defer file.Close()
//
//	// Вызываем ExternalSort
//	err = ExternalSort(file, os.Stdout, *k, *n, *r, *u)
//	if err != nil {
//		log.Fatalf("Error sorting: %v", err)
//	}
//}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	filename  = "test.txt"
	chunkSize = 1000
)

// GetColumn - получаем колонку из строки
func GetColumn(line string, n int) string {
	cols := strings.Split(line, "\t")
	if n-1 < len(cols) {
		return cols[n-1]
	}

	return ""
}

// SortLines - сортировка среза строк в памяти
func SortLines(lines []string, column int, numeric, reverse bool) {
	sort.Slice(lines, func(i, j int) bool {
		col1 := GetColumn(lines[i], column)
		col2 := GetColumn(lines[j], column)

		if numeric {
			ai, errA := strconv.Atoi(col1)
			bi, errB := strconv.Atoi(col2)
			if errA == nil && errB == nil {
				if reverse {
					return ai > bi
				}
				return ai < bi
			}
		}

		if reverse {
			return col1 > col2
		}

		return col1 < col2
	})
}

// ReadChunk - чтение куска строк из сканера
func ReadChunk(scanner *bufio.Scanner, chunkSize int) ([]string, bool) {
	lines := make([]string, 0, chunkSize)
	for i := 0; i < chunkSize && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	return lines, len(lines) > 0
}

// ExternalSort - сортировка кусками
func ExternalSort(filename string, column int, numeric, reverse, unique bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var chunkFiles []string

	// создаём временные отсортированные файлы
	for {
		lines, ok := ReadChunk(scanner, chunkSize)
		if !ok {
			break
		}

		SortLines(lines, column, numeric, reverse)

		tmpFile, err := os.CreateTemp("", "chunk_*.txt")
		if err != nil {
			return err
		}

		for _, line := range lines {
			fmt.Fprintln(tmpFile, line)
		}
		tmpFile.Close()
		chunkFiles = append(chunkFiles, tmpFile.Name())
	}

	// слияние временных файлов
	type fileLine struct {
		line string
		scan *bufio.Scanner
		file *os.File
	}

	openFiles := []fileLine{}
	for _, fname := range chunkFiles {
		f, _ := os.Open(fname)
		sc := bufio.NewScanner(f)
		if sc.Scan() {
			openFiles = append(openFiles, fileLine{line: sc.Text(), scan: sc, file: f})
		} else {
			f.Close()
		}
	}

	var lastPrinted string
	for len(openFiles) > 0 {
		// находим минимальную или максимальную строку
		minIdx := 0
		for i := 1; i < len(openFiles); i++ {
			a := GetColumn(openFiles[i].line, column)
			b := GetColumn(openFiles[minIdx].line, column)
			var cmp bool
			if numeric {
				ai, errA := strconv.Atoi(a)
				bi, errB := strconv.Atoi(b)
				if errA == nil && errB == nil {
					cmp = ai < bi
				} else {
					cmp = a < b
				}
			} else {
				cmp = a < b
			}
			if reverse {
				cmp = !cmp
			}
			if cmp {
				minIdx = i
			}
		}

		// печатаем укникальные строки
		if !unique || openFiles[minIdx].line != lastPrinted {
			fmt.Println(openFiles[minIdx].line)
			lastPrinted = openFiles[minIdx].line
		}

		// читаем следующую строку из того же файла
		if openFiles[minIdx].scan.Scan() {
			openFiles[minIdx].line = openFiles[minIdx].scan.Text()
		} else {
			openFiles[minIdx].file.Close()
			openFiles = append(openFiles[:minIdx], openFiles[minIdx+1:]...)
		}
	}

	// удаляем временные файлы
	for _, fname := range chunkFiles {
		_ = os.Remove(fname)
	}

	return nil
}

func main() {
	// Определяем флаги
	k := flag.Int("k", 1, "sort by column number")
	n := flag.Bool("n", false, "numeric sort")
	r := flag.Bool("r", false, "reverse sort")
	u := flag.Bool("u", false, "unique lines only")

	// комбинированные флаги типа -nr
	flag.BoolVar(n, "nr", false, "numeric reverse sort ")
	flag.BoolVar(r, "rn", false, "reverse numeric sort ")
	flag.BoolVar(n, "nru", false, "numeric reverse unique sort ")

	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Вызываем ExternalSort
	err = ExternalSort(filename, *k, *n, *r, *u)
	if err != nil {
		log.Fatalf("Error sorting: %v", err)
	}
}
