package entity

import (
	"time"

	"github.com/google/uuid"
)

// Notify - структура данных для уведомлений
type Notify struct {
	UserID      string `json:"user_id"`
	NotifyID    uuid.UUID
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Status      bool      `json:"status"`
	SendingDate time.Time `json:"sending_date"`
	RetryCount  uint8     `json:"retry_count"`
}

// NotifyCache - структура для кеша с расширенными полями
type NotifyCache struct {
	UserID     string    `json:"user_id"`
	NotifyID   uuid.UUID `json:"notify_id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Status     bool      `json:"status"`
	EventTime  time.Time `json:"event_time"`
	RetryCount uint8     `json:"retry_count"`
}
