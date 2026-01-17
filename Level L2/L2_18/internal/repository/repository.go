package repository

import (
	"L2_18/internal/entity"
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

	return nil
}

// GetEventForMonth - получить событие на месяц
func (repo *Repository) GetEventForMonth(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	// проверяем что такой id есть
	userEvents, ok := repo.Storage[userID]
	if !ok {
		return nil, entity.ErrNoEvents
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
		return nil, entity.ErrNoEvents
	}

	return result, nil
}

// GetEventForDay - получаем событие на день
func (repo *Repository) GetEventForDay(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	userEvents, ok := repo.Storage[userID]
	if !ok {
		return nil, entity.ErrNoEvents
	}

	// нормализуем дату для поиска
	normalizedDate := date.Truncate(24 * time.Hour)

	// прямой поиск по ключу вместо перебора
	events, ok := userEvents[normalizedDate]
	if !ok || len(events) == 0 {
		return nil, entity.ErrNoEvents
	}

	return events, nil
}

// GetEventsForWeek - получаем событие на неделю
func (repo *Repository) GetEventsForWeek(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	userEvents, ok := repo.Storage[userID]
	if !ok {
		return nil, entity.ErrNoEvents
	}

	var result []entity.Calendar

	// находим понедельник недели
	weekday := int(date.Weekday())
	if weekday == 0 { // если воскресенье
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
		return nil, entity.ErrNoEvents
	}

	return result, nil
}

// UpdateEvent - обновить событие
func (repo *Repository) UpdateEvent(event entity.Calendar) error {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	userEvents, ok := repo.Storage[event.UserID]
	if !ok {
		return entity.ErrNoEvents
	}

	data := event.DataEvent.Truncate(24 * time.Hour)

	events, ok := userEvents[data]
	if !ok || len(events) == 0 {
		return entity.ErrNoEvents
	}

	// ищем событие для обновления
	found := false
	for i := range events {
		if events[i].NameEvent == event.NameEvent && events[i].DataEvent.Equal(event.DataEvent) {
			events[i].Text = event.Text
			// если нужно обновить имя, то обновляем
			if event.NameEvent != "" {
				events[i].NameEvent = event.NameEvent
			}
			found = true
			break
		}
	}

	if !found {
		return entity.ErrEventNotFound
	}

	// сохраняем то что обновили
	userEvents[data] = events

	return nil
}

// DeleteEvent - удалить событие
func (repo *Repository) DeleteEvent(event entity.Calendar) error {
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	userEvents, ok := repo.Storage[event.UserID]
	if !ok {
		return entity.ErrNoEvents
	}

	data := event.DataEvent.Truncate(24 * time.Hour)

	events, ok := userEvents[data]
	if !ok || len(events) == 0 {
		return entity.ErrNoEvents
	}

	// убираем из среза не нужное событие(проверяем имя и дату)
	filtered := make([]entity.Calendar, 0, len(events))
	deleted := false
	for _, e := range events {
		if e.NameEvent == event.NameEvent && e.DataEvent.Equal(event.DataEvent) {
			deleted = true
			continue // пропускаем это событие
		}
		filtered = append(filtered, e)
	}

	if !deleted {
		return entity.ErrEventNotFound
	}

	// если больше нет событий на дату удаляем ключ
	if len(filtered) == 0 {
		delete(userEvents, data)
	} else {
		// сохраняем измененный срез
		userEvents[data] = filtered
	}

	return nil
}
