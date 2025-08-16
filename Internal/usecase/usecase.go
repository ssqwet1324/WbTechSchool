package usecase

import (
	"WbDemoProject/Internal/entity"
	"context"
	"fmt"
)

type RepositoryProvider interface {
	SaveOrderInDB(ctx context.Context, order *entity.Order) error
	GetOrderFromCache(orderUID string) (*entity.Order, bool)
	GetOrderFromDB(ctx context.Context, orderUID string) (*entity.Order, error)
}
type Usecase struct {
	repo RepositoryProvider
}

func New(repo RepositoryProvider) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

// SaveOrderInDB - сохраняем данные
func (u *Usecase) SaveOrderInDB(ctx context.Context, order *entity.Order) error {
	if err := u.repo.SaveOrderInDB(ctx, order); err != nil {
		return fmt.Errorf("failed to save order in DB: %w", err)
	}

	return nil
}

// GetOrderFromCache - получаем данные из кэша
func (u *Usecase) GetOrderFromCache(orderUID string) (*entity.Order, bool) {
	data, exist := u.repo.GetOrderFromCache(orderUID)
	if !exist {
		return nil, false
	}

	return data, true
}

// GetOrderFromDB - получаем заказ из дб
func (u *Usecase) GetOrderFromDB(ctx context.Context, orderUID string) (*entity.Order, error) {
	data, err := u.repo.GetOrderFromDB(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order from DB: %w", err)
	}

	return data, nil
}

// HandleOrder - имплементирующем интерфейс из кафки с заказом
func (u *Usecase) HandleOrder(ctx context.Context, order *entity.Order) error {
	//тут проверить на валидные данные
	return u.SaveOrderInDB(ctx, order)
}
