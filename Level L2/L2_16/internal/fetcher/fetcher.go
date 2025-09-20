package fetcher

import (
	"errors"
	"io"
	"log"
	"net/http"
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
func (fetcher *Fetcher) Fetch(url string) ([]byte, error) {
	req, err := fetcher.client.Get(url)
	if err != nil {
		return nil, errors.New("Fetch: error get request: " + url + ": " + err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Fetch: Error closing body", err)
		}
	}(req.Body)

	if req.StatusCode != http.StatusOK {
		return nil, errors.New("Fetch: error fetching " + url + ": " + req.Status)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New("Fetch: error reading body: " + url + ": " + err.Error())
	}

	return body, nil
}
