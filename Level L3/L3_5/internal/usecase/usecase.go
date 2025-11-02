package usecase

import (
	"context"
	"errors"
	"event_booker/internal/entity"
	"fmt"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)

type RepositoryProvider interface {
	CreateEvent(ctx context.Context, event entity.CreateEvent, totalNumberSeats int, seatNumbers []int) error
	CheckFreeSeats(ctx context.Context, eventID string, seatNumber int) (string, error)
	MarkSeatAsReserving(ctx context.Context, eventID string, seatNumber int, userID string) error
	ConfirmSeatBooking(ctx context.Context, eventID string, seatNumber int, userID string) error
	CleanupExpiredReservations(ctx context.Context, eventID string, seatNumber int, redisExist bool) error
	GetEventInfo(ctx context.Context, eventID string) (*entity.EventInfo, error)
	GetAllEvents(ctx context.Context) ([]entity.CreateEvent, error)

	AddReserveSeat(ctx context.Context, key string, seat entity.Seat) (bool, error)
	CheckReserveUser(ctx context.Context, key string, seat entity.Seat) error
	RemoveReserveSeat(ctx context.Context, key string) error
	CheckKeyInRedis(ctx context.Context, key string) (bool, error)
}

type UseCase struct {
	repo RepositoryProvider
}

func New(repo RepositoryProvider) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// generateSeatNumbers — создаёт список номеров мест
func generateSeatNumbers(rows, seatsPerRow, startNumber int) []int {
	var seats []int
	for r := 0; r < rows; r++ {
		rowStart := startNumber + r*seatsPerRow
		for s := 0; s < seatsPerRow; s++ {
			seats = append(seats, rowStart+s)
		}
	}
	return seats
}

// generateRedisKey — генерирует ключ для Redis
func generateRedisKey(eventID string, seatNumber int) string {
	return fmt.Sprintf("seat_lock:%s:%d", eventID, seatNumber)
}

// CreateEvent - создать мероприятие
func (uc *UseCase) CreateEvent(ctx context.Context, event entity.CreateEvent, layout entity.TotalSeats) (string, error) {
	eventID := uuid.New().String()
	event.ID = eventID

	// генерация мест
	seats := generateSeatNumbers(layout.Rows, layout.SeatsPerRow, layout.StartNumber)

	// общее количество мест
	totalNumberSeats := layout.Rows * layout.SeatsPerRow

	// создаем мероприятие
	err := uc.repo.CreateEvent(ctx, event, totalNumberSeats, seats)
	if err != nil {
		return "", err
	}

	zlog.Logger.Info().Msgf("Event created: ID=%s, Title=%s, Date=%s", event.ID, event.Title, event.Date)

	return event.ID, nil
}

// TryReserveSeat - забронировать место
func (uc *UseCase) TryReserveSeat(ctx context.Context, eventID string, seat entity.Seat) error {
	key := generateRedisKey(eventID, seat.SeatNumber)

	// пробуем заблокировать в Redis
	locked, err := uc.repo.AddReserveSeat(ctx, key, seat)
	if err != nil {
		// Технический сбой Redis
		return fmt.Errorf("service unavailable: %w", err)
	}

	if !locked {
		// Место уже бронируется
		return errors.New("seat is already being reserved by another user")
	}

	// если ошибка удаляем блокировку
	defer func() {
		if err != nil {
			// Пытаемся освободить, но не перезаписываем оригинальную ошибку
			if releaseErr := uc.repo.RemoveReserveSeat(ctx, key); releaseErr != nil {
				zlog.Logger.Error().Err(releaseErr).Str("key", key).Msg("failed to release seat lock")
			}
		}
	}()

	// Проверяем БД на свободное место
	seatStatus, err := uc.repo.CheckFreeSeats(ctx, eventID, seat.SeatNumber)
	if err != nil {
		return fmt.Errorf("failed to check seat status: %w", err)
	}

	if seatStatus != "free" {
		return errors.New("seat already booked")
	}

	// Обновляем статус в БД с проверкой
	err = uc.repo.MarkSeatAsReserving(ctx, eventID, seat.SeatNumber, seat.UserID)
	if err != nil {
		return fmt.Errorf("failed to reserve seat in DB: %w", err)
	}

	return nil
}

// ConfirmSeatBooking - подтверждаем бронирование места
func (uc *UseCase) ConfirmSeatBooking(ctx context.Context, eventID string, seatNumber int, userID string) error {
	// Проверяем, что бронь принадлежит пользователю в Redis
	key := generateRedisKey(eventID, seatNumber)
	err := uc.repo.CheckReserveUser(ctx, key, entity.Seat{
		SeatNumber: seatNumber,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("reservation not found or expired: %w", err)
	}

	// Подтверждаем бронирование в БД
	err = uc.repo.ConfirmSeatBooking(ctx, eventID, seatNumber, userID)
	if err != nil {
		return err
	}

	// Удаляем временную блокировку из Redis
	if releaseErr := uc.repo.RemoveReserveSeat(ctx, key); releaseErr != nil {
		zlog.Logger.Error().Err(releaseErr).Str("key", key).Msg("failed to release seat lock after confirmation")
	}

	return nil
}

// GetEvent - получение информации о событии и свободных местах
func (uc *UseCase) GetEvent(ctx context.Context, eventID string) (*entity.EventInfo, error) {
	return uc.repo.GetEventInfo(ctx, eventID)
}

// GetAllEvents - получаем все мероприятия
func (uc *UseCase) GetAllEvents(ctx context.Context) ([]entity.CreateEvent, error) {
	return uc.repo.GetAllEvents(ctx)
}

func (uc *UseCase) CleanupExpiredReservations(ctx context.Context, eventID string, seatNumber int) error {
	key := generateRedisKey(eventID, seatNumber)

	// проверяем, есть ли ключ в Redis
	existKey, err := uc.repo.CheckKeyInRedis(ctx, key)
	if err != nil {
		return fmt.Errorf("ошибка проверки ключа в Redis: %w", err)
	}

	// если ключ есть — просто выходим (бронь ещё активна)
	if existKey {
		return nil
	}

	// если ключа нет — очищаем конкретное место
	err = uc.repo.CleanupExpiredReservations(ctx, eventID, seatNumber, existKey)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("event_id", eventID).Int("seat_number", seatNumber).Msg("Ошибка очистки просроченного места")

		return fmt.Errorf("ошибка очистки просроченного места: %w", err)
	}

	return nil
}
