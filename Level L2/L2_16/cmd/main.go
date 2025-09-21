package main

import (
	"L2_16/internal/crawler"
	"L2_16/internal/fetcher"
	"L2_16/internal/reader"
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	// чтобы не падать делаем в бесконечном цикле
	for {
		url, depth, err := reader.ReadCommand()
		if err != nil {
			log.Println("Ошибка чтения:", err)
			continue
		}
		depthInt, err := strconv.Atoi(depth)
		if err != nil {
			log.Fatal("download depth is not a number", err)
		}

		fmt.Printf("Начинаем скачивание с глубиной %d: %s\n", depthInt, url)

		// создаем состояние краулера
		state := crawler.NewCrawlerState(url)

		// создаем f
		f := fetcher.New(30 * time.Second)

		// запускаем рекурсивное скачивание
		err = crawler.CrawlPage(url, depthInt, state, f)
		if err != nil {
			log.Printf("Ошибка скачивания: %v", err)
			continue
		}

		fmt.Printf("Скачивание завершено! Обработано %d страниц\n", state.VisitedCount())
	}
}
