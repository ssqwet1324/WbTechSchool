package migrations

import (
	"context"
	"fmt"
	"shortener/internal/config"
	"shortener/internal/repository"
	"time"

	"github.com/wb-go/wbf/zlog"
)

// Migration - структура миграций
type Migration struct {
	repo *repository.Repository
	cfg  *config.ServiceConfig
}

// New - конструктор
func New(repo *repository.Repository, config *config.ServiceConfig) *Migration {
	return &Migration{repo: repo, cfg: config}
}

// InitTables - создаёт таблицы для shortener
func (m *Migration) InitTables(ctx context.Context) error {
	zlog.Logger.Info().Msg("Начинаем создание таблиц для shortener")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS short_urls (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS clicks (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(255) NOT NULL REFERENCES short_urls(short_url) ON DELETE CASCADE,
			user_agent TEXT,
			created_at TIMESTAMP DEFAULT now()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_short_url ON clicks(short_url);`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_created_at ON clicks(created_at);`,
		`CREATE TABLE IF NOT EXISTS clicks_aggregate (
			short_url VARCHAR(255) PRIMARY KEY,
			total_clicks INT DEFAULT 0,
			clicks_by_day JSONB DEFAULT '{}'::jsonb,
			clicks_by_month JSONB DEFAULT '{}'::jsonb,
			clicks_by_user_agent JSONB DEFAULT '{}'::jsonb
		);`,
	}

	maxRetries := m.cfg.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3 // Дефолтное значение
	}

	for i := 0; i < maxRetries; i++ {
		success := true
		for _, q := range queries {
			if _, err := m.repo.DB.ExecContext(ctx, q); err != nil {
				zlog.Logger.Error().Err(err).Int("attempt", i+1).
					Int("max_attempts", maxRetries).
					Msg("Ошибка выполнения миграции")
				success = false
				break
			}
		}
		if success {
			zlog.Logger.Warn().Msg("Таблицы short_urls, clicks и clicks_aggregate успешно созданы или уже существуют")
			return nil
		}
		time.Sleep(m.cfg.RetryDelay)
	}

	return fmt.Errorf("не удалось создать таблицы после %d попыток", maxRetries)
}
