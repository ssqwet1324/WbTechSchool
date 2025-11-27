package worker

import (
	"L4_3/internal/entity"
	"L4_3/internal/log"
	"time"
)

// UseCase - интерфейс для работы с функциями usecase
type UseCase interface {
	MoveOldEventsToArchive(before time.Time) error
}

// Worker - структура воркера
type Worker struct {
	uc            UseCase
	logger        *log.Log
	RemainderChan chan entity.Event
}

// New - конструктор воркера
func New(logger *log.Log, uc UseCase, ch chan entity.Event) *Worker {
	return &Worker{
		uc:            uc,
		logger:        logger,
		RemainderChan: ch,
	}
}

// RunReminderWorker - запуск напоминания
func (w *Worker) RunReminderWorker() {
	go func() {
		for event := range w.RemainderChan {
			go func(ev entity.Event) {
				now := time.Now()
				if ev.RemindAt.After(now) {
					time.Sleep(time.Until(ev.RemindAt))
				}

				if w.logger != nil {
					w.logger.AsyncMessagef("[REMINDER] user=%s event=%s text=%s remind_at=%s",
						ev.UserID,
						ev.NameEvent,
						ev.Text,
						ev.RemindAt.Format(time.RFC3339),
					)
				}
			}(event)
		}
	}()
}

// RunCleanerWorker - берет старые события и помещает в архив
func (w *Worker) RunCleanerWorker(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := w.uc.MoveOldEventsToArchive(time.Now()); err != nil && w.logger != nil {
				w.logger.AsyncError("could not move old events to archive", err)
			}
		}
	}()
}
