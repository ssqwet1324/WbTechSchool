package migrations

import (
	"context"
	"fmt"
	"image_processor/internal/config"
	"image_processor/internal/repository"
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

// InitTables - создаёт таблицы для image_processor
func (m *Migration) InitTables(ctx context.Context) error {
	zlog.Logger.Info().Msg("Начинаем создание таблиц для обработки изображений")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS photos (
    	id VARCHAR(255) PRIMARY KEY,
    	name TEXT NOT NULL,
    	status VARCHAR(50) NOT NULL DEFAULT 'waiting',
    	created_at TIMESTAMP DEFAULT now()
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
