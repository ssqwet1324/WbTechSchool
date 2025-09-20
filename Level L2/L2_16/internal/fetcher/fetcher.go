package fetcher

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Fetcher - запрос клиент-сервер
type Fetcher struct {
	client *http.Client
}

// New - конструктор Fetcher
func New(timeout time.Duration) *Fetcher {
	client := &http.Client{
		Timeout:   timeout,
		Transport: http.DefaultTransport,
	}
	return &Fetcher{
		client: client,
	}
}

// Fetch - получаем тело страницы (body)
func (fetcher *Fetcher) Fetch(link string) ([]byte, string, error) {
	req, err := fetcher.client.Get(link)
	if err != nil {
		return nil, "", errors.New("Fetch: error get request: " + link + ": " + err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Fetch: Error closing body", err)
		}
	}(req.Body)

	// проверяем что смогли подключиться
	if req.StatusCode != http.StatusOK {
		return nil, "", errors.New("Fetch: error fetching " + link + ": " + req.Status)
	}

	// получаем body сайта
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, "", errors.New("Fetch: error reading body: " + link + ": " + err.Error())
	}

	// парсим ссылку для получения домена
	u, err := url.Parse(link)
	if err != nil {
		return nil, "", errors.New("Fetch: error parsing domain: " + link + ": " + err.Error())
	}

	// получаем домен
	domain := u.Scheme + "://" + u.Host

	return body, domain, nil
}
