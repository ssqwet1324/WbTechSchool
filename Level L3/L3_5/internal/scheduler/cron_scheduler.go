package scheduler

import (
	"context"
	"event_booker/internal/usecase"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/wb-go/wbf/zlog"
)

// Scheduler — отвечает за очистку просроченных броней
type Scheduler struct {
	uc   *usecase.UseCase
	cron *cron.Cron
}

// New создаёт новый планировщик
func New(uc *usecase.UseCase) *Scheduler {
	return &Scheduler{
		uc:   uc,
		cron: cron.New(),
	}
}

// Start запускает cron, который каждые 5 минут чистит просроченные брони
func (s *Scheduler) Start() {
	_, err := s.cron.AddFunc("@every 5m", func() {
		ctx := context.Background()
		zlog.Logger.Info().Msg("[Scheduler] Запуск проверки всех мероприятий...")

		events, err := s.uc.GetAllEvents(ctx)
		if err != nil {
			zlog.Logger.Error().Msgf("Ошибка получения мероприятий: %v", err)
			return
		}

		// создаем горутины для каждого мероприятия
		var wg sync.WaitGroup
		for _, event := range events {
			wg.Add(1)
			go func(eID, title string, total int) {
				defer wg.Done()

				zlog.Logger.Info().Msgf("Проверка мероприятия: %s (%s)", title, eID)

				for seat := 1; seat <= total; seat++ {
					err := s.uc.CleanupExpiredReservations(ctx, eID, seat)
					if err != nil {
						zlog.Logger.Error().Msgf("Ошибка очистки (event=%s seat=%d): %v", eID, seat, err)
					}
				}

				zlog.Logger.Info().Msgf("Очистка завершена для мероприятия: %s", title)

			}(event.ID, event.Title, event.TotalSeats)
		}

		wg.Wait()
		zlog.Logger.Info().Msg("[Scheduler] Очистка всех мероприятий завершена.")
	})

	if err != nil {
		zlog.Logger.Fatal().Msgf("Ошибка при добавлении cron-задачи: %v", err)
	}

	s.cron.Start()
	zlog.Logger.Info().Msg(" Планировщик запущен (каждые 5 минут)")
}

// Stop — остановка cron-планировщика
func (s *Scheduler) Stop() {
	s.cron.Stop()
	zlog.Logger.Info().Msg("Планировщик остановлен")
}
