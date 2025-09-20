package main

import (
	"L2_16/internal/fetcher"
	"L2_16/internal/parser"
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

		body, domain, err := fetcher.New(30 * time.Second).Fetch(url)
		if err != nil {
			log.Println("Ошибка загрузки:", err)
			continue
		}

		// тут формируем файлы всех страниц по глубине рекурсии
		if depthInt > 1 {
			res, err := parser.Parser(body, domain)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, link := range res.Link {
				fmt.Println(link)
			}

			fmt.Println(res.CSS)
		} else {
			fmt.Println(string(body))
			// тут мы будем формировать файл только той страницы на которую указан url
		}

		fmt.Println("домен:", domain)
	}
}
