package entity

import "time"

// NewItem — используется для создания новой записи
type NewItem struct {
	Title    string    `json:"title"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
	Category string    `json:"category"`
}

// Item — для хранения в бд
type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"` // Дата и время операции
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"` // Дата и время создания записи
	UpdatedAt time.Time `json:"updated_at"` // Время обновления
}

// GetItems  - получение записей от 1й даты до 2й
type GetItems struct {
	FromDate *time.Time `json:"from_date"`
	ToDate   *time.Time `json:"to_date"`
}

// GetItemsFromAnalytics - получение аналитики от 1й даты до 2й + по категории
type GetItemsFromAnalytics struct {
	FromDate *time.Time `json:"from_date"`
	ToDate   *time.Time `json:"to_date"`
	Category *string    `json:"category"` // как дополнение
}

// AnalyticsResult - результат работы аналитики
type AnalyticsResult struct {
	TotalCount int64
	TotalSum   float64
	AvgAmount  float64
	Median     float64
	P90        float64
}
