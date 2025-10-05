package migrations

import (
	"L3_1/internal/config"
	"L3_1/internal/repository"
	"context"
	"fmt"
	"log"
	"time"
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

// InitNotifyTable создаёт таблицу notifications
func (m *Migration) InitNotifyTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS notifications (
		user_id TEXT NOT NULL,
		notify_id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT,
		status BOOLEAN NOT NULL DEFAULT FALSE,
		sending_date TIMESTAMP NOT NULL,
		retry_count SMALLINT NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_notifications_sending_date ON notifications (sending_date);`

	maxRetries := m.cfg.MaxRetries
	retryDelay := m.cfg.RetryDelay * time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = m.repo.DB.ExecContext(ctx, query)
		if err == nil {
			log.Println("Таблица notifications успешно создана или уже существует")
			return nil
		}

		log.Printf("Ошибка создания таблицы notifications (попытка %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("не удалось создать таблицу notifications после %d попыток: %w", maxRetries, err)
}
