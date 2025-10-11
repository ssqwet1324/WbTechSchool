package entity

// ShortenURL - структура для ссылок
type ShortenURL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// ShortenURLAnalytics - структура для аналитики
type ShortenURLAnalytics struct {
	ShortURL    string         `json:"short_url"`
	TotalClicks int            `json:"total_clicks"`
	ByDay       map[string]int `json:"by_day"`
	ByMonth     map[string]int `json:"by_month"`
	ByUserAgent map[string]int `json:"by_user_agent"`
}
