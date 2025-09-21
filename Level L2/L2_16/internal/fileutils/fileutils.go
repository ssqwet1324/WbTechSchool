package fileutils

import (
	"net/url"
	"os"
	"path"
	"strings"
)

// CreateFilePath создает путь для сохранения файла на основе URL
func CreateFilePath(urlStr, resourceType string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return path.Join("mirror", resourceType, "error.html")
	}

	// создаем путь на основе URL
	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")

	// если путь пустой, используем index.html
	if len(pathParts) == 1 && pathParts[0] == "" {
		pathParts = []string{"index.html"}
	} else {
		// добавляем расширение к последней части, если нет расширения
		lastPart := pathParts[len(pathParts)-1]
		if !strings.Contains(lastPart, ".") {
			// для HTML файлов добавляем .html, для остальных - по умолчанию
			if resourceType == "pages" {
				pathParts[len(pathParts)-1] = lastPart + ".html"
			}
		}
	}

	// добавляем query параметры если есть
	if u.RawQuery != "" {
		lastIdx := len(pathParts) - 1
		pathParts[lastIdx] = pathParts[lastIdx] + "_" + url.QueryEscape(u.RawQuery)
	}

	return path.Join("mirror", resourceType, path.Join(pathParts...))
}

// CreateFileName создает имя файла из URL
func CreateFileName(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return "index.html"
	}

	filename := path.Base(u.Path)
	if filename == "" || filename == "/" {
		filename = "index.html"
	}

	// если нет расширения — добавляем .html
	if !strings.Contains(filename, ".") {
		filename += ".html"
	}

	// учитываем query параметры
	if u.RawQuery != "" {
		filename += "_" + url.QueryEscape(u.RawQuery)
	}

	return filename
}

// SaveFile сохраняет файл с созданием необходимых директорий
func SaveFile(body []byte, filePath string) error {
	// создаем директории если нужно
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	// сохраняем файл
	return os.WriteFile(filePath, body, 0644)
}
