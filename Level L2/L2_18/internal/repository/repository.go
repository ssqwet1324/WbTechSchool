package repository

import (
	"L2_18/internal/entity"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Repository - структура для хранилища
type Repository struct {
	Storage map[string]map[time.Time][]entity.Calendar
	Mutex   *sync.RWMutex
	Log     *zap.Logger
}

// New - конструктор структуры
func New(log *zap.Logger) *Repository {
	return &Repository{
		Storage: make(map[string]map[time.Time][]entity.Calendar),
		Mutex:   &sync.RWMutex{},
		Log:     log.Named("Repository"),
	}
}

// SaveEvent - сохраняем событие в календаре
func (repo *Repository) SaveEvent(event entity.Calendar) error {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	// проверяем есть ли хранилище для пользователя
	if _, ok := repo.Storage[event.UserID]; !ok {
		repo.Storage[event.UserID] = make(map[time.Time][]entity.Calendar)
	}

	// парсим дату корректно
	data := event.DataEvent.Truncate(24 * time.Hour)

	// сохраняем по userID событие и дату
	repo.Storage[event.UserID][data] = append(repo.Storage[event.UserID][data], event)

	repo.Log.Info("Saving event", zap.String("event_name", event.NameEvent))

	return nil
}

// GetEventForMonth - получить событие на месяц
func (repo *Repository) GetEventForMonth(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	// проверяем что такой id есть
	userEvents, ok := repo.Storage[userID]
	if !ok {
		return nil, errors.New("user not found")
	}

	var result []entity.Calendar
	year, month, _ := date.Date()

	// проходимся по срезу и сравниваем дату
	for day, events := range userEvents {
		y, m, _ := day.Date()
		if y == year && m == month {
			result = append(result, events...)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no events for this month")
	}

	return result, nil
}

// GetEventForDay - получаем событие на день
func (repo *Repository) GetEventForDay(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	userEvents, ok := repo.Storage[userID]
	if !ok {
		return nil, errors.New("user not found")
	}

	var result []entity.Calendar
	year, month, day := date.Date()

	// проходимся по срезу и сравниваем дату
	for eventDay, events := range userEvents {
		y, m, d := eventDay.Date()
		if y == year && m == month && d == day {
			result = append(result, events...)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no events for this day")
	}

	return result, nil
}

// GetEventsForWeek - получаем событие на неделю
func (repo *Repository) GetEventsForWeek(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	userEvents, ok := repo.Storage[userID]
	if !ok {
		return nil, errors.New("user not found")
	}

	var result []entity.Calendar

	// находим понедельник недели
	weekday := int(date.Weekday())
	if weekday == 0 { // если воскресенье, в Go это 0
		weekday = 7
	}

	startOfWeek := date.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 6) // воскресенье

	for eventDay, events := range userEvents {
		day := eventDay.Truncate(24 * time.Hour)
		if !day.Before(startOfWeek) && !day.After(endOfWeek) {
			result = append(result, events...)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no events for this week")
	}

	return result, nil
}

// UpdateEvent - обновить событие
func (repo *Repository) UpdateEvent(event entity.Calendar) error {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	userEvents, ok := repo.Storage[event.UserID]
	if !ok {
		return errors.New("user not found")
	}

	data := event.DataEvent.Truncate(24 * time.Hour)

	events, ok := userEvents[data]
	if !ok {
		return errors.New("no events on this date")
	}

	// тут обновляем данные у события
	for i := range events {
		if events[i].DataEvent == event.DataEvent {
			events[i].NameEvent = event.NameEvent
			events[i].Text = event.Text
			break
		}
	}

	// сохраняем то что обновили
	userEvents[data] = events

	repo.Log.Info("Updating event", zap.String("event_name", event.NameEvent))

	return nil
}

// DeleteEvent - удалить событие
func (repo *Repository) DeleteEvent(event entity.Calendar) error {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	userEvents, ok := repo.Storage[event.UserID]
	if !ok {
		return errors.New("user not found")
	}

	data := event.DataEvent.Truncate(24 * time.Hour)

	events, ok := userEvents[data]
	if !ok {
		return errors.New("no events on this date")
	}

	fmt.Println("список events", events)

	// убираем из среза не нужное событие(проверяем имя и дату)
	filtered := make([]entity.Calendar, 0, len(events))
	for _, e := range events {
		if !(e.NameEvent == event.NameEvent && e.DataEvent.Equal(event.DataEvent)) {
			filtered = append(filtered, e)
		}
	}

	// если больше нет событий на дату удаляем ключ
	if len(filtered) == 0 {
		delete(userEvents, data)
	} else {
		// сохраняем измененный срез
		userEvents[data] = filtered
	}

	repo.Log.Info("Deleting event", zap.String("event", event.NameEvent), zap.String("date", data.String()))

	return nil
}
