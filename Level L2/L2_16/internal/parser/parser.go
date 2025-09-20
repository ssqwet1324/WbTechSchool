package parser

import (
	"bytes"
	"errors"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// Parser - парсим все ссылки на странице
func Parser(body []byte) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("Parse: error parsing body: " + err.Error())
	}

	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if !strings.HasPrefix(a.Val, "http://") && !strings.HasPrefix(a.Val, "https://") {
						newLink := ResolveURL(a.Val)
						links = append(links, newLink)
					} else {
						links = append(links, a.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	urls := DeleteDuplies(links)

	return urls, nil
}

// ResolveURL - меняем относительный url на полный
func ResolveURL(relative string) string {
	base, _ := url.Parse("https://skillbox.ru")
	rel, err := url.Parse(relative)
	if err != nil {
		return relative // если парсинг не удался, возвращаем как есть
	}
	return base.ResolveReference(rel).String()
}

// DeleteDuplies - удаляем дубликаты из среза
func DeleteDuplies(urls []string) []string {
	slices.Sort(urls)
	sliceWithoutDuplies := slices.Compact(urls)

	return sliceWithoutDuplies
}
