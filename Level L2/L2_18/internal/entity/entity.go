package entity

import (
	"errors"
	"time"
)

// Calendar - структура хранения данных
type Calendar struct {
	UserID    string    `json:"user_id"`
	NameEvent string    `json:"name_event"`
	DataEvent time.Time `json:"data_event"`
	Text      string    `json:"text"`
}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrNoEvents      = errors.New("no events")
	ErrParsing       = errors.New("parsing error")
)
