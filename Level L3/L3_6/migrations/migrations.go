package migrations

import (
	"context"
	"fmt"
	"sales_tracker/internal/config"
	"sales_tracker/internal/repository"
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
		`CREATE TABLE IF NOT EXISTS items (
        id UUID PRIMARY KEY,
        title TEXT NOT NULL,
        amount NUMERIC NOT NULL CHECK (amount >= 0),
        date TIMESTAMP NOT NULL,
        category TEXT,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
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
