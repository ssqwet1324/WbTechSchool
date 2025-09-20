package parser

import (
	"bytes"
	"errors"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// Resources - структура списков для хранения элементов внутри тегов
type Resources struct {
	CSS  []string
	JS   []string
	Img  []string
	Link []string
}

// Parser - парсим все ресурсы на странице
func Parser(body []byte, domain string) (*Resources, error) {
	resources := &Resources{
		CSS:  make([]string, 0),
		JS:   make([]string, 0),
		Img:  make([]string, 0),
		Link: make([]string, 0),
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("Parse: error parsing body: " + err.Error())
	}

	var links []string

	// для рекурсивного обхода
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// находим все основные теги для файлов
			switch n.Data {
			case "a":
				for _, a := range n.Attr {
					if a.Key == "href" {
						if !strings.HasPrefix(a.Val, "http://") && !strings.HasPrefix(a.Val, "https://") {
							links = append(links, resolveURL(a.Val, domain))
						} else {
							links = append(links, a.Val)
						}
					}
				}
			case "script":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resources.JS = append(resources.JS, resolveURL(attr.Val, domain))
					}
				}
			case "img":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resources.Img = append(resources.Img, resolveURL(attr.Val, domain))
					}
				}
			case "link":
				var href, rel string
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
					}
					if attr.Key == "rel" {
						rel = attr.Val
					}
				}
				if rel == "stylesheet" {
					resources.CSS = append(resources.CSS, resolveURL(href, domain))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// удаляем дубликаты ссылок
	resources.Link = deleteDuplies(links)

	return resources, nil
}

// resolveURL - создаем из относительной ссылки полную
func resolveURL(relative, domain string) string {
	base, err := url.Parse(domain)
	if err != nil {
		return relative
	}
	rel, err := url.Parse(relative)
	if err != nil {
		return relative
	}

	return base.ResolveReference(rel).String()
}

// deleteDuplies - удаляем дубликаты url
func deleteDuplies(urls []string) []string {
	slices.Sort(urls)

	return slices.Compact(urls)
}
