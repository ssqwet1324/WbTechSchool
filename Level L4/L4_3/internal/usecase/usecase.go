package usecase

import (
	"L4_3/internal/entity"
	"fmt"
	"time"
)

// RepositoryProvider - интерфейс для работы с хранилищем событий календаря
type RepositoryProvider interface {
	SaveEvent(event entity.Calendar) error
	UpdateEvent(event entity.Calendar) error
	DeleteEvent(event entity.Calendar) error
	GetEventForDay(userID string, date time.Time) ([]entity.Calendar, error)
	GetEventsForWeek(userID string, date time.Time) ([]entity.Calendar, error)
	GetEventForMonth(userID string, date time.Time) ([]entity.Calendar, error)
	AddOldEventInArchive(event entity.Calendar) error
	MoveOldEventsToArchive(before time.Time) error
}

// UseCase - структура сервиса
type UseCase struct {
	provider      RepositoryProvider
	remainderChan chan entity.Event
}

// New - конструктор для UseCase
func New(provider RepositoryProvider, ch chan entity.Event) *UseCase {
	return &UseCase{
		provider:      provider,
		remainderChan: ch,
	}
}

// SaveEvent - сохранить событие
func (uc *UseCase) SaveEvent(event entity.Calendar) error {
	err := uc.provider.SaveEvent(event)
	if err != nil {
		return fmt.Errorf("could not save event: %w", err)
	}

	if !event.RemindAt.IsZero() {
		uc.remainderChan <- entity.Event{
			UserID:    event.UserID,
			NameEvent: event.NameEvent,
			Text:      event.Text,
			RemindAt:  event.RemindAt,
		}
	}

	return nil
}

// UpdateEvent - обновить событие
func (uc *UseCase) UpdateEvent(event entity.Calendar) error {
	err := uc.provider.UpdateEvent(event)
	if err != nil {
		return fmt.Errorf("there is no such event: %w", err)
	}

	return nil
}

// DeleteEvent - удалить событие
func (uc *UseCase) DeleteEvent(event entity.Calendar) error {
	err := uc.provider.DeleteEvent(event)
	if err != nil {
		return fmt.Errorf("there is no such event: %w", err)
	}

	return nil
}

// GetEventForDay - получить события на день
func (uc *UseCase) GetEventForDay(userID, date string) ([]entity.Calendar, error) {
	dateTime, err := ParseDate(date)
	if err != nil {
		return nil, fmt.Errorf("could not parse data: %w", err)
	}

	events, err := uc.provider.GetEventForDay(userID, dateTime)
	if err != nil {
		return nil, fmt.Errorf("could not get event for day: %w", err)
	}

	return events, nil
}

// GetEventsForWeek - получить события на неделю
func (uc *UseCase) GetEventsForWeek(userID, date string) ([]entity.Calendar, error) {
	dateTime, err := ParseDate(date)
	if err != nil {
		return nil, fmt.Errorf("could not parse data: %w", err)
	}

	events, err := uc.provider.GetEventsForWeek(userID, dateTime)
	if err != nil {
		return nil, fmt.Errorf("could not get events for week: %w", err)
	}

	return events, nil
}

// GetEventForMonth - получить события на месяц
func (uc *UseCase) GetEventForMonth(userID, date string) ([]entity.Calendar, error) {
	dateTime, err := ParseDate(date)
	if err != nil {
		return nil, fmt.Errorf("could not parse data: %w", err)
	}

	events, err := uc.provider.GetEventForMonth(userID, dateTime)
	if err != nil {
		return nil, fmt.Errorf("could not get event for month: %w", err)
	}

	return events, nil
}

// ParseDate - парсим дату
func ParseDate(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, fmt.Errorf("date is empty")
	}
	if len(date) != 10 {
		return time.Time{}, fmt.Errorf("date is invalid")
	}

	return time.Parse("2006-01-02", date)
}

// AddOldEventInArchive - убираем в архив старое событие
func (uc *UseCase) AddOldEventInArchive(event entity.Calendar) error {
	return uc.provider.AddOldEventInArchive(event)
}

// MoveOldEventsToArchive - переносим старые события в архив
func (uc *UseCase) MoveOldEventsToArchive(before time.Time) error {
	return uc.provider.MoveOldEventsToArchive(before)
}
