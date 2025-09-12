package main

import (
	"L2_12/internal/cli"
	"L2_12/internal/search"
	"bufio"
	"log"
	"os"
)

func main() {
	options, pattern, filename := cli.ParseOptions()

	var scanner *bufio.Scanner

	if filename == "" || filename == "-" {
		// Если файла нет или указан "-", читаем из stdin
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		// Иначе открываем файл
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("Ошибка открытия файла:", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	search.PatternString(scanner, &options, pattern)
}
