package main

import (
	"bufio"
	"log"
	"os"

	"L4.2/internal/app"
	"L4.2/internal/cli"
)

func main() {
	flag, pattern, filename, err := cli.ParseServerFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Парсим grep флаги (они уже объявлены в ParseServerFlags через flag.*)
	options := cli.ParseOptions()

	switch flag.Mode {
	case "worker":
		if err := app.RunWorker(flag); err != nil {
			log.Fatal(err)
		}
		return
	case "leader":
	default:
		log.Fatalf("неизвестный режим %q", flag.Mode)
	}

	var scanner *bufio.Scanner
	if filename == "" || filename == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Ошибка открытия файла %s: %v", filename, err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	if err := app.RunServer(flag, options, pattern, scanner); err != nil {
		log.Fatal(err)
	}
}
