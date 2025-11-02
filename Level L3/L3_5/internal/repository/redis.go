package repository

import (
	"context"
	"errors"
	"event_booker/internal/entity"
	"time"

	"github.com/wb-go/wbf/zlog"
)

// AddReserveSeat - добавляем временно бронирования место до подтверждения оплаты
// Возвращаем (успех bool, ошибка error)
func (repo *Repository) AddReserveSeat(ctx context.Context, key string, seat entity.Seat) (bool, error) {
	//TODO В конфиг
	ok, err := repo.Client.SetNX(ctx, key, seat.UserID, 10*time.Minute).Result()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("redis add reserve seat error")
		return false, err
	}

	if !ok {
		// место уже забронировано
		zlog.Logger.Warn().Str("key", key).Msg("seat is temporarily reserved by another user")
		return false, nil
	}

	return true, nil
}

// CheckReserveUser - проверка на то, кто бронирует место
func (repo *Repository) CheckReserveUser(ctx context.Context, key string, seat entity.Seat) error {
	// Проверяем что бронь принадлежит этому пользователю
	lockedBy, _ := repo.Client.Get(ctx, key)
	if lockedBy != seat.UserID {
		zlog.Logger.Error().Msg("reservation expired or belongs to another user")

		return errors.New("reservation expired or belongs to another user")
	}

	return nil
}

// CheckKeyInRedis - проверяем есть ли такой ключ в redis
func (repo *Repository) CheckKeyInRedis(ctx context.Context, key string) (bool, error) {
	ok, err := repo.Client.Exists(ctx, key).Result()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("redis check exists error")
		return false, err
	}
	if ok == 0 {
		return false, nil
	}

	return true, nil
}

// RemoveReserveSeat - удаляем бронь
func (repo *Repository) RemoveReserveSeat(ctx context.Context, key string) error {
	return repo.Client.Del(ctx, key)
}
