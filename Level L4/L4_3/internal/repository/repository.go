package repository

import (
	"L4_3/internal/entity"
	"L4_3/internal/log"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Repository - структура для хранилища
type Repository struct {
	storage map[string]map[time.Time][]entity.Calendar
	archive map[string]map[time.Time][]entity.Calendar
	mutex   *sync.RWMutex
	log     *log.Log
}

// New - конструктор структуры
func New(log *log.Log) *Repository {
	return &Repository{
		storage: make(map[string]map[time.Time][]entity.Calendar),
		archive: make(map[string]map[time.Time][]entity.Calendar),
		mutex:   &sync.RWMutex{},
		log:     log,
	}
}

// SaveEvent - сохраняем событие в календаре
func (repo *Repository) SaveEvent(event entity.Calendar) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// проверяем есть ли хранилище для пользователя
	if _, ok := repo.storage[event.UserID]; !ok {
		repo.storage[event.UserID] = make(map[time.Time][]entity.Calendar)
	}

	// парсим дату корректно
	data := event.DataEvent.Truncate(24 * time.Hour)

	// сохраняем по userID событие и дату
	repo.storage[event.UserID][data] = append(repo.storage[event.UserID][data], event)

	repo.log.AsyncMessage("Repository: Save Event: " + event.UserID)

	return nil
}

// GetEventForMonth - получить событие на месяц
func (repo *Repository) GetEventForMonth(userID string, date time.Time) ([]entity.Calendar, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	// проверяем что такой id есть
	userEvents, ok := repo.storage[userID]
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
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	userEvents, ok := repo.storage[userID]
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
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	userEvents, ok := repo.storage[userID]
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
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	userEvents, ok := repo.storage[event.UserID]
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

	repo.log.AsyncMessage("Repository: Update Event: " + event.UserID)

	return nil
}

// DeleteEvent - удалить событие
func (repo *Repository) DeleteEvent(event entity.Calendar) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	userEvents, ok := repo.storage[event.UserID]
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

	repo.log.AsyncMessage("Repository: Delete Event: " + event.UserID + "date" + data.String())

	return nil
}

// AddOldEventInArchive - добавляем новый старый event в архив
func (repo *Repository) AddOldEventInArchive(event entity.Calendar) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.addEventToArchiveLocked(event)

	if err := repo.checkTotalArchiveLocked(event.UserID); err != nil {
		return errors.New("error checking total archive")
	}

	repo.log.AsyncMessage("Repository: Add Old Event: " + event.NameEvent)

	return nil
}

// MoveOldEventsToArchive - перемещает все события старше before в архив
func (repo *Repository) MoveOldEventsToArchive(before time.Time) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	cutoff := before.Truncate(24 * time.Hour)

	for userID, userEvents := range repo.storage {
		for day, events := range userEvents {
			dayTruncated := day.Truncate(24 * time.Hour)
			if dayTruncated.Before(cutoff) {
				for _, event := range events {
					repo.addEventToArchiveLocked(event)
				}
				delete(userEvents, day)
			}
		}

		if len(userEvents) == 0 {
			delete(repo.storage, userID)
		}

		if err := repo.checkTotalArchiveLocked(userID); err != nil {
			return errors.New("error checking total archive")
		}
	}

	return nil
}

func (repo *Repository) addEventToArchiveLocked(event entity.Calendar) {
	// проверяем есть ли архив для пользователя
	if _, ok := repo.archive[event.UserID]; !ok {
		repo.archive[event.UserID] = make(map[time.Time][]entity.Calendar)
	}

	data := event.DataEvent.Truncate(24 * time.Hour)

	repo.archive[event.UserID][data] = append(repo.archive[event.UserID][data], event)
}

// checkTotalArchiveLocked - проверяем количество элементов архива; вызываем под блокировкой
func (repo *Repository) checkTotalArchiveLocked(userID string) error {
	userArchive := repo.archive[userID]
	if len(userArchive) <= 10 {
		return nil
	}

	// находим самый старый день
	var oldest time.Time
	first := true
	for date := range userArchive {
		if first || date.Before(oldest) {
			oldest = date
			first = false
		}
	}

	// удаляем самый старый день
	delete(userArchive, oldest)

	repo.log.AsyncMessage("Repository: Deleted oldest archive day" + "user" + userID + "deleted_day" + oldest.String())

	return nil
}
