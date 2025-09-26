package usecase

import (
	"L2_18/internal/entity"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// RepositoryProvider - интерфейс для работы с хранилищем событий календаря
type RepositoryProvider interface {
	SaveEvent(event entity.Calendar) error
	UpdateEvent(event entity.Calendar) error
	DeleteEvent(event entity.Calendar) error
	GetEventForDay(userID string, date time.Time) ([]entity.Calendar, error)
	GetEventsForWeek(userID string, date time.Time) ([]entity.Calendar, error)
	GetEventForMonth(userID string, date time.Time) ([]entity.Calendar, error)
}

// UseCase - структура сервиса
type UseCase struct {
	provider RepositoryProvider
	log      *zap.Logger
}

// New - конструктор для UseCase
func New(provider RepositoryProvider, logger *zap.Logger) *UseCase {
	return &UseCase{
		provider: provider,
		log:      logger,
	}
}

// SaveEvent - сохранить событие
func (uc *UseCase) SaveEvent(event entity.Calendar) error {
	err := uc.provider.SaveEvent(event)
	if err != nil {
		return fmt.Errorf("could not save event: %w", err)
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

// GetEventsForWeek - получить события на недел
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
