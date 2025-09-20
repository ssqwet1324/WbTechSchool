package main

import (
	"L2_16/internal/fetcher"
	"L2_16/internal/parser"
	"L2_16/internal/reader"
	"fmt"
	"log"
	"time"
)

func main() {
	// чтобы не падать делаем в бесконечном цикле
	for {
		url, err := reader.ReadCommand()
		if err != nil {
			log.Println("Ошибка чтения:", err)
			continue
		}

		body, err := fetcher.New(30 * time.Second).Fetch(url)
		if err != nil {
			log.Println("Ошибка загрузки:", err)
			continue
		}

		links, err := parser.Parser(body)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, link := range links {
			fmt.Println(link)
		}
	}
}
