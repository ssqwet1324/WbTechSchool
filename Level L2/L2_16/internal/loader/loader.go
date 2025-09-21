package loader

import (
	"L2_16/internal/fetcher"
	"L2_16/internal/fileutils"
	"L2_16/internal/parser"
	"L2_16/internal/queue"
	"fmt"
	"sync"
	"time"
)

// DownloadFile - загружаем все ссылки из очереди параллельно
func DownloadFile(queue <-chan string, resourceType string) {
	var wg sync.WaitGroup
	workers := 5
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for link := range queue {
				// создаём путь с подпапкой
				filePath := fileutils.CreateFilePath(link, resourceType)

				fmt.Printf("[Worker %d] Downloading %s -> %s\n", workerID, link, filePath)

				// скачиваем данные с url
				downloadBody, _, err := fetcher.New(30 * time.Second).Fetch(link)
				if err != nil {
					fmt.Println("Error downloading", link, err)
					continue
				}

				// сохраняем файл
				err = fileutils.SaveFile(downloadBody, filePath)
				if err != nil {
					fmt.Println("Error saving", link, err)
					continue
				}

				fmt.Printf("[Worker %d] Saved: %s\n", workerID, filePath)
			}
		}(i)
	}
	wg.Wait()
}

// DownloadPageResources скачивает ресурсы страницы (CSS, JS, изображения)
func DownloadPageResources(res *parser.Resources, baseDomain string) error {
	// скачиваем CSS
	if len(res.CSS) > 0 {
		cssQueue, err := queue.Queue(res.CSS)
		if err == nil {
			go DownloadFile(cssQueue, "css")
		}
	}

	// скачиваем JS
	if len(res.JS) > 0 {
		jsQueue, err := queue.Queue(res.JS)
		if err == nil {
			go DownloadFile(jsQueue, "js")
		}
	}

	// скачиваем изображения
	if len(res.Img) > 0 {
		imgQueue, err := queue.Queue(res.Img)
		if err == nil {
			go DownloadFile(imgQueue, "img")
		}
	}

	return nil
}
