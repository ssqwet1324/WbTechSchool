package crawler

import (
	"L2_16/internal/fetcher"
	"L2_16/internal/fileutils"
	"L2_16/internal/loader"
	"L2_16/internal/parser"
	"L2_16/internal/queue"
	"fmt"
	"log"
	"net/url"
	"sync"
)

// State тут храним краулер
type State struct {
	visited    map[string]bool
	mu         sync.RWMutex
	baseDomain string
}

// NewCrawlerState инициализируем краулер
func NewCrawlerState(baseURL string) *State {
	u, _ := url.Parse(baseURL)
	return &State{
		visited:    make(map[string]bool),
		baseDomain: u.Host,
	}
}

// IsVisited - проверяет, был ли URL уже посещен
func (cs *State) IsVisited(url string) bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.visited[url]
}

// MarkVisited - помечает URL как посещенный
func (cs *State) MarkVisited(url string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.visited[url] = true
}

// IsSameDomain - проверяет, принадлежит ли URL тому же домену
func (cs *State) IsSameDomain(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return u.Host == cs.baseDomain
}

// VisitedCount - возвращает количество посещенных страниц
func (cs *State) VisitedCount() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return len(cs.visited)
}

// CrawlPage - рекурсивно скачивает страницу и все связанные ресурсы
func CrawlPage(urlStr string, depth int, state *State, fetcher *fetcher.Fetcher) error {
	// проверяем, не посещали ли мы уже эту страницу
	if state.IsVisited(urlStr) {
		return nil
	}

	// проверяем, что URL принадлежит тому же домену
	if !state.IsSameDomain(urlStr) {
		return nil
	}

	// помечаем URL как посещенный
	state.MarkVisited(urlStr)

	fmt.Printf("Скачиваем страницу (глубина %d): %s\n", depth, urlStr)

	// загружаем страницу
	body, domain, err := fetcher.Fetch(urlStr)
	if err != nil {
		return fmt.Errorf("ошибка загрузки %s: %v", urlStr, err)
	}

	// парсим страницу
	includeLinks := depth > 1
	res, err := parser.Parser(body, domain, includeLinks)
	if err != nil {
		return fmt.Errorf("ошибка парсинга %s: %v", urlStr, err)
	}

	// создаем путь для сохранения HTML файла
	htmlPath := fileutils.CreateFilePath(urlStr, "pages")

	// сохраняем HTML файл
	err = fileutils.SaveFile(body, htmlPath)
	if err != nil {
		return fmt.Errorf("ошибка сохранения HTML %s: %v", urlStr, err)
	}

	fmt.Printf("✓ Сохранен HTML: %s\n", htmlPath)

	// скачиваем ресурсы страницы (CSS, JS, изображения)
	err = loader.DownloadPageResources(res, state.baseDomain)
	if err != nil {
		log.Printf("Ошибка скачивания ресурсов для %s: %v", urlStr, err)
	}

	// если глубина больше 1, рекурсивно скачиваем найденные ссылки
	if depth > 1 {
		// создаем очередь для страниц (исключая текущую)
		var pageLinks []string
		for _, link := range res.Link {
			if link != urlStr { // пропускаем текущую страницу
				pageLinks = append(pageLinks, link)
			}
		}

		if len(pageLinks) > 0 {
			pagesQueue, err := queue.Queue(pageLinks)
			if err == nil {
				go loader.DownloadFile(pagesQueue, "pages")
			}
		}
	}

	return nil
}
