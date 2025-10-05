package worker

import (
	"L3_1/internal/entity"
	"L3_1/internal/rabbit"
	"L3_1/internal/usecase"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)

// Worker - структура воркера
type Worker struct {
	ctx       context.Context
	useCase   *usecase.UseCaseNotify
	publisher *rabbit.Publisher
	wakeup    chan bool
}

// New - конструктор
func New(ctx context.Context, useCase *usecase.UseCaseNotify, publisher *rabbit.Publisher) *Worker {
	return &Worker{
		ctx:       ctx,
		useCase:   useCase,
		publisher: publisher,
		wakeup:    make(chan bool, 1),
	}
}

// Run - воркер
func (w *Worker) Run() {
	for {
		select {
		case <-w.ctx.Done():
			zlog.Logger.Info().Msg("Worker exit")
			return
		default:
			// получаем ближайшее уведомление
			notify, err := w.getNextNotify()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			now := time.Now()
			if notify.EventTime.Before(now) || notify.EventTime.Equal(now) {
				// если время уже пришло или прошло — отправляем сразу
				zlog.Logger.Info().
					Str("notifyID", notify.NotifyID.String()).
					Time("eventTime", notify.EventTime).
					Time("now", now).
					Msg("EventTime уже прошёл, отправляем сразу")
				w.SendNotification(*notify)
				continue
			}

			// ждем, пока наступит время уведомления или сигнал wakeup
			duration := time.Until(notify.EventTime)
			timer := time.NewTimer(duration)
			select {
			case <-w.ctx.Done():
				timer.Stop()
				return
			case <-w.wakeup:
				zlog.Logger.Info().Msg("Worker wakeup signal received")
				timer.Stop()
				// перейти к следующей итерации, чтобы заново выбрать ближайшее уведомление
				continue
			case <-timer.C:
				// время пришло — перед отправкой убеждаемся, что уведомление всё ещё актуально
				w.SendNotification(*notify)
			}
		}
	}
}

// getNextNotify - получаем следующие уведомление
func (w *Worker) getNextNotify() (*entity.NotifyCache, error) {
	return w.useCase.GetNearNotify(w.ctx)
}

// SendNotification - отправляем уведомление
func (w *Worker) SendNotification(notify entity.NotifyCache) {
	// проверка на актуальность в кеше уведомления
	cached, err := w.useCase.GetNotifyInCash(w.ctx, notify.NotifyID.String())
	if err == nil {
		// Если нет в кэше — значит удален
		if cached.NotifyID == uuid.Nil {
			zlog.Logger.Info().Str("notifyID", notify.NotifyID.String()).Msg("Уведомление отсутствует в кэше, отправка отменена")
			return
		}
		// Если время изменилось на будущее — не отправляем сейчас
		if cached.EventTime.After(time.Now()) {
			zlog.Logger.Info().
				Str("notifyID", notify.NotifyID.String()).
				Time("newEventTime", cached.EventTime).
				Msg("Уведомление переотложено, отправка перенесена")
			return
		}
	}
	data, err := json.Marshal(map[string]interface{}{
		"user_id":    notify.UserID,
		"id":         notify.NotifyID.String(),
		"title":      notify.Title,
		"body":       notify.Body,
		"event_time": notify.EventTime,
	})

	// повторная отправка
	if err != nil {
		zlog.Logger.Err(err).Str("title", notify.Title).Msg("Ошибка маршала уведомления")
		w.RetryNotify(notify)
		return
	}

	// отправляем в очередь
	if err := w.publisher.Publish("notify_queue", data); err != nil {
		zlog.Logger.Err(err).Str("title", notify.Title).Msg("Ошибка отправки в RabbitMQ")
		w.RetryNotify(notify)
		return
	}

	// удаляем из бд
	if err := w.useCase.DeleteNotification(w.ctx, notify.NotifyID.String()); err != nil {
		zlog.Logger.Err(err).Str("title", notify.Title).Msg("Не удалось удалить уведомление из БД после отправки")
	}

	// удаляем из кеша
	if err := w.useCase.DeleteNotifyInCash(w.ctx, notify.NotifyID.String()); err != nil {
		zlog.Logger.Err(err).Str("title", notify.Title).Msg("Не удалось удалить уведомление из кэша после отправки")
	}

	zlog.Logger.Info().Str("title", notify.Title).Msg("Уведомление отправлено")
}

// RetryNotify - повторяем отправку
func (w *Worker) RetryNotify(notify entity.NotifyCache) {
	if notify.RetryCount >= 5 {
		zlog.Logger.Warn().Str("title", notify.Title).Msg("Превышено количество попыток, уведомление не будет отправлено")
		return
	}

	notify.RetryCount++
	delay := time.Duration(1<<notify.RetryCount) * time.Second
	if delay > 5*time.Minute {
		delay = 5 * time.Minute
	}
	notify.EventTime = time.Now().Add(delay)

	if err := w.useCase.AddNotifyInCash(w.ctx, notify); err != nil {
		zlog.Logger.Err(err).Str("title", notify.Title).Msg("Не удалось обновить уведомление в кэше для повторной отправки")
	}
}

// WakeUpWorker - пробуждает воркер
func (w *Worker) WakeUpWorker() {
	select {
	case w.wakeup <- true:
		zlog.Logger.Info().Msg("Сигнал на пробуждение воркера отправлен")
	default:
		zlog.Logger.Warn().Msg("Сигнал на пробуждение воркера уже отправлен")
	}
}
