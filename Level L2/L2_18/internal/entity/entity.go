package entity

import "time"

// Calendar - структура хранения данных
type Calendar struct {
	UserID    string    `json:"user_id"`
	NameEvent string    `json:"name_event"`
	Text      string    `json:"text"`
	DataEvent time.Time `json:"data_event"`
}
