package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"shortener/internal/entity"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

// Repository - структура для работы с БД
type Repository struct {
	DB     *dbpg.DB
	master *dbpg.Options
}

// New - конструктор для repository
func New(masterDSN string, options *dbpg.Options) *Repository {
	masterDB, err := sql.Open("postgres", masterDSN)
	if err != nil {
		log.Fatalf("failed to open master DB: %v", err)
	}

	masterDB.SetMaxOpenConns(options.MaxOpenConns)
	masterDB.SetMaxIdleConns(options.MaxIdleConns)
	masterDB.SetConnMaxLifetime(options.ConnMaxLifetime)

	db := &dbpg.DB{
		Master: masterDB,
		Slaves: nil,
	}

	return &Repository{
		DB:     db,
		master: options,
	}
}

// AddShortUrl - добавить ссылку в бд
func (repo *Repository) AddShortUrl(ctx context.Context, urls entity.ShortenURL) error {
	query := `INSERT INTO short_urls (original_url, short_url) VALUES ($1, $2)`
	res, err := repo.DB.ExecContext(ctx, query, urls.OriginalURL, urls.ShortURL)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: AddShortUrl: failed to add short url")
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Repository: AddShortUrl: failed to check affected rows")
		return err
	}
	if rows == 0 {
		zlog.Logger.Warn().Msg("Repository: AddShortUrl: no rows affected")
		return err
	}

	return nil
}

// GetShortUrl - получить короткий url по оригинальному
func (repo *Repository) GetShortUrl(ctx context.Context, originalURL string) (string, error) {
	query := `SELECT short_url FROM short_urls WHERE original_url=$1`
	var shortURL string
	err := repo.DB.QueryRowContext(ctx, query, originalURL).Scan(&shortURL)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			zlog.Logger.Error().Err(err).Str("original_url", originalURL).Msg("Repository: GetShortUrl: failed to get short url")
		}
		return "", err
	}

	return shortURL, nil
}

// GetOriginalURL - получить оригинальный url по короткому
func (repo *Repository) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	query := `SELECT original_url FROM short_urls WHERE short_url=$1`
	var originalURL string
	err := repo.DB.QueryRowContext(ctx, query, shortURL).Scan(&originalURL)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			zlog.Logger.Error().Err(err).Str("short_url", shortURL).Msg("Repository: GetOriginalURL: failed to get original url")
		}
		return "", err
	}

	return originalURL, nil
}

// ExistsShortUrl - проверить есть ли такой url в бд
func (repo *Repository) ExistsShortUrl(ctx context.Context, shortUrl string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM short_urls WHERE short_url = $1)`
	var exists bool
	err := repo.DB.QueryRowContext(ctx, query, shortUrl).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// AddAnalytics - добавление аналитики для короткой ссылки
func (repo *Repository) AddAnalytics(ctx context.Context, urls entity.ShortenURLAnalytics) error {
	// Получаем существующую запись
	var agg struct {
		Total       int
		ByDay       []byte
		ByMonth     []byte
		ByUserAgent []byte
	}

	err := repo.DB.QueryRowContext(ctx, `
		SELECT total_clicks, clicks_by_day, clicks_by_month, clicks_by_user_agent
		FROM clicks_aggregate WHERE short_url = $1`, urls.ShortURL).
		Scan(&agg.Total, &agg.ByDay, &agg.ByMonth, &agg.ByUserAgent)

	existingByDay := map[string]int{}
	existingByMonth := map[string]int{}
	existingByUserAgent := map[string]int{}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		zlog.Logger.Error().Err(err).Str("short_url", urls.ShortURL).Msg("select failed")
		return err
	}

	// если есть данные — распаковываем JSON
	_ = json.Unmarshal(agg.ByDay, &existingByDay)
	_ = json.Unmarshal(agg.ByMonth, &existingByMonth)
	_ = json.Unmarshal(agg.ByUserAgent, &existingByUserAgent)

	// соединяем новые значения
	agg.Total += urls.TotalClicks
	for k, v := range urls.ByDay {
		existingByDay[k] += v
	}
	for k, v := range urls.ByMonth {
		existingByMonth[k] += v
	}
	for k, v := range urls.ByUserAgent {
		existingByUserAgent[k] += v
	}

	// Маршалим обратно в JSON
	byDayJSON, _ := json.Marshal(existingByDay)
	byMonthJSON, _ := json.Marshal(existingByMonth)
	byUAJSON, _ := json.Marshal(existingByUserAgent)

	// Вставка или обновление
	if errors.Is(err, sql.ErrNoRows) {
		_, err = repo.DB.ExecContext(ctx, `
			INSERT INTO clicks_aggregate
			(short_url, total_clicks, clicks_by_day, clicks_by_month, clicks_by_user_agent)
			VALUES ($1, $2, $3::jsonb, $4::jsonb, $5::jsonb)
		`, urls.ShortURL, agg.Total, byDayJSON, byMonthJSON, byUAJSON)
	} else {
		_, err = repo.DB.ExecContext(ctx, `
			UPDATE clicks_aggregate
			SET total_clicks = $2, clicks_by_day = $3::jsonb, clicks_by_month = $4::jsonb, clicks_by_user_agent = $5::jsonb
			WHERE short_url = $1
		`, urls.ShortURL, agg.Total, byDayJSON, byMonthJSON, byUAJSON)
	}

	if err != nil {
		zlog.Logger.Error().Err(err).Str("short_url", urls.ShortURL).Msg("write failed")
	}

	return err
}

// GetAnalytics - получить аналитику о ссылке
func (repo *Repository) GetAnalytics(ctx context.Context, shortUrl string) (*entity.ShortenURLAnalytics, error) {
	var agg struct {
		Total       int
		ByDay       []byte
		ByMonth     []byte
		ByUserAgent []byte
	}

	query := `SELECT total_clicks, clicks_by_day, clicks_by_month, clicks_by_user_agent FROM clicks_aggregate WHERE short_url = $1`

	err := repo.DB.QueryRowContext(ctx, query, shortUrl).Scan(
		&agg.Total,
		&agg.ByDay,
		&agg.ByMonth,
		&agg.ByUserAgent,
	)
	if err != nil {
		return nil, err
	}

	// Распаковка JSONB в map
	byDay := map[string]int{}
	byMonth := map[string]int{}
	byUserAgent := map[string]int{}

	_ = json.Unmarshal(agg.ByDay, &byDay)
	_ = json.Unmarshal(agg.ByMonth, &byMonth)
	_ = json.Unmarshal(agg.ByUserAgent, &byUserAgent)

	analytics := &entity.ShortenURLAnalytics{
		ShortURL:    shortUrl,
		TotalClicks: agg.Total,
		ByDay:       byDay,
		ByMonth:     byMonth,
		ByUserAgent: byUserAgent,
	}

	return analytics, nil
}
