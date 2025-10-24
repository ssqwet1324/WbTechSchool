package entity

import (
	"io"
	"time"
)

// LoadPhoto - структура запроса для загрузки фото
type LoadPhoto struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Reader    io.Reader
}

// PhotoInfo - структура для получения информации о фото
type PhotoInfo struct {
	BucketName string `json:"bucket_name"`
	Version    string `json:"version"`
	PhotoID    string `json:"photo_id"`
	Width      string `json:"width"`
	Height     string `json:"height"`
}

// MinIOObject - объект для формирования url
type MinIOObject struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}
