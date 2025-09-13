package main

import (
	"L2_13/internal/analoguecut"
	"L2_13/internal/cli"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	flags, fileName := cli.ParseFlags()

	var scanner *bufio.Scanner
	if fileName == "" || fileName == "-" {
		// Если файла нет или указан "-", читаем из stdin
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		// Иначе открываем файл
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal("Ошибка открытия файла:", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	str := analoguecut.Cut(scanner, flags)
	fmt.Println(str)
}
