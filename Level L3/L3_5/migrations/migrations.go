package migrations

import (
	"context"
	"event_booker/internal/config"
	"event_booker/internal/repository"
	"fmt"
	"time"

	"github.com/wb-go/wbf/zlog"
)

// Migration - структура миграций
type Migration struct {
	repo *repository.Repository
	cfg  *config.Config
}

// New - конструктор
func New(repo *repository.Repository, config *config.Config) *Migration {
	return &Migration{repo: repo, cfg: config}
}

// InitTables - создаёт таблицы для image_processor
func (m *Migration) InitTables(ctx context.Context) error {
	zlog.Logger.Info().Msg("Начинаем создание таблиц для обработки изображений")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		date TIMESTAMP NOT NULL,
		total_seats INT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
	);`,

		`CREATE TABLE IF NOT EXISTS seats (
		id TEXT PRIMARY KEY,
		event_id TEXT REFERENCES events(id) ON DELETE CASCADE,
		seat_number INT NOT NULL,
		status TEXT NOT NULL DEFAULT 'free',
		user_id TEXT NOT NULL,
		UNIQUE (event_id, seat_number)
	);`,
	}

	maxRetries := m.cfg.MaxRetries
	if maxRetries == 0 {
		maxRetries = 5 // Дефолтное значение
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
			zlog.Logger.Warn().Msg("Таблица photos успешно создана или уже существует")
			return nil
		}

		time.Sleep(m.cfg.RetryDelay)
	}

	return fmt.Errorf("не удалось создать таблицы после %d попыток", maxRetries)
}
