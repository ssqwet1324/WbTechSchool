package migrations

import (
	"comment_tree/internal/config"
	"comment_tree/internal/repository"
	"context"
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

// InitTables - создаёт таблицы для shortener
func (m *Migration) InitTables(ctx context.Context) error {
	zlog.Logger.Info().Msg("Начинаем создание таблиц для комментариев")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS flat_comments (
    	id VARCHAR(255) PRIMARY KEY,
    	text TEXT NOT NULL,
    	created_at TIMESTAMP DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
    	id VARCHAR(255) PRIMARY KEY,
    	parent_id VARCHAR(255) REFERENCES comments(id) ON DELETE CASCADE,
    	comment_ref VARCHAR(255) NOT NULL REFERENCES flat_comments(id) ON DELETE CASCADE,
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
			zlog.Logger.Warn().Msg("Таблицы flat_comments, comments успешно созданы или уже существуют")
			return nil
		}

		time.Sleep(m.cfg.RetryDelay)
	}

	return fmt.Errorf("не удалось создать таблицы после %d попыток", maxRetries)
}
