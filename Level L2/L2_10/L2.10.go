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
	k := flag.Int("k", 1, "search by column number")
	n := flag.Bool("n", false, "numeric search")
	r := flag.Bool("r", false, "reverse search")
	u := flag.Bool("u", false, "unique lines only")

	// комбинированные флаги типа -nr
	flag.BoolVar(n, "nr", false, "numeric reverse search ")
	flag.BoolVar(r, "rn", false, "reverse numeric search ")
	flag.BoolVar(n, "nru", false, "numeric reverse unique search ")

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
