package parser

import (
	"bytes"
	"errors"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// Resources хранит списки всех найденных ресурсов на странице
type Resources struct {
	CSS  []string
	JS   []string
	Img  []string
	Link []string
}

// Parser - парсит HTML и собирает все ресурсы (CSS, JS, img, ссылки)
func Parser(body []byte, domain string, includeLinks bool) (*Resources, error) {
	resources := &Resources{
		CSS:  make([]string, 0),
		JS:   make([]string, 0),
		Img:  make([]string, 0),
		Link: make([]string, 0),
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("Parser: error parsing body: " + err.Error())
	}

	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				if includeLinks {
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							if !strings.HasPrefix(attr.Val, "http://") && !strings.HasPrefix(attr.Val, "https://") {
								links = append(links, resolveURL(attr.Val, domain))
							} else {
								links = append(links, attr.Val)
							}
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
				for _, attr := range n.Attr {
					if attr.Key == "rel" && attr.Val == "stylesheet" {
						for _, a := range n.Attr {
							if a.Key == "href" {
								resources.CSS = append(resources.CSS, resolveURL(a.Val, domain))
							}
						}
					} else if includeLinks && attr.Key == "rel" && (attr.Val == "canonical" || attr.Val == "alternate") {
						for _, a := range n.Attr {
							if a.Key == "href" {
								links = append(links, resolveURL(a.Val, domain))
							}
						}
					}
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

// resolveURL создаёт абсолютный URL из относительного
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

// deleteDuplies удаляет дубли в списке ссылок
func deleteDuplies(urls []string) []string {
	slices.Sort(urls)
	return slices.Compact(urls)
}
