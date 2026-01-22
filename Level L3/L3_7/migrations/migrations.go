package migrations

import (
	"fmt"
	"time"
	"warehouse_control/internal/config"
	"warehouse_control/internal/repository"

	"github.com/pressly/goose/v3"
	"github.com/wb-go/wbf/zlog"
)

// Migration - миграции
type Migration struct {
	repo *repository.Repository
	cfg  *config.Config
}

// New - конструктор миграций
func New(repo *repository.Repository, cfg *config.Config) *Migration {
	return &Migration{
		repo: repo,
		cfg:  cfg,
	}
}

// RunMigrations выполняет все миграции из папки
func (m *Migration) RunMigrations() error {
	// путь к папке с .sql файлами
	dir := m.cfg.PackageWithMigrations

	maxRetries := m.cfg.MaxRetries
	if maxRetries == 0 {
		maxRetries = 5
	}

	for i := 0; i < maxRetries; i++ {
		if err := goose.Up(m.repo.DB.Master, dir); err != nil {
			zlog.Logger.Error().Err(err).Int("attempt", i+1).Int("max_attempts", maxRetries).
				Msg("Ошибка выполнения миграций")
			time.Sleep(m.cfg.RetryDelay)
			continue
		}

		zlog.Logger.Info().Msg("Все миграции успешно применены")
		return nil
	}

	return fmt.Errorf("не удалось применить миграции после %d попыток", maxRetries)
}
