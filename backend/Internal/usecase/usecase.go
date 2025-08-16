package usecase

import (
	"WbDemoProject/Internal/entity"
	"context"
	"fmt"
	"log"
)

type RepositoryProvider interface {
	SaveOrderInDB(ctx context.Context, order *entity.Order) error
	GetOrderFromCache(orderUID string) (*entity.Order, bool)
	GetOrderFromDB(ctx context.Context, orderUID string) (*entity.Order, error)
	SaveOrderInCache(order *entity.Order) error
}

type Usecase struct {
	repo RepositoryProvider
}

func New(repo RepositoryProvider) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

// SaveOrderInDB - сохраняем данные заказа в БД
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

// CheckOrderFromCacheInDB - проверяем есть ли заказ который в кэше в БД
// ВАЖНО: защищает от утечек данных - если заказ есть в кэше, но нет в БД
func (u *Usecase) CheckOrderFromCacheInDB(ctx context.Context, orderUID string, cacheResponse bool) error {
	// Если в кэше заказа нет, то выходим
	if !cacheResponse {
		return nil
	}

	// Получаем данные из кэша
	data, exist := u.repo.GetOrderFromCache(orderUID)
	if !exist {
		return nil
	}

	// Проверяем БД
	_, err := u.GetOrderFromDB(ctx, orderUID)
	if err != nil {
		// Если в БД нет записи - добавляем (защита от утечек!)
		go func(order *entity.Order) {
			if err := u.repo.SaveOrderInDB(context.Background(), order); err != nil {
				log.Printf("Failed to save order from cache to DB: %v", err)
			}
		}(data)
	}

	return nil
}

// SaveOrderInCache - сохраняем заказ в кэше
func (u *Usecase) SaveOrderInCache(order *entity.Order) error {
	if err := u.repo.SaveOrderInCache(order); err != nil {
		return fmt.Errorf("failed to save order in cache: %w", err)
	}

	return nil
}

// GetOrderFromDB - получаем заказ из БД и автоматически кэшируем
func (u *Usecase) GetOrderFromDB(ctx context.Context, orderUID string) (*entity.Order, error) {
	data, err := u.repo.GetOrderFromDB(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order from DB: %w", err)
	}

	// сохраняем в кэш
	if err = u.SaveOrderInCache(data); err != nil {
		return nil, err
	}

	return data, nil
}

// GetOrder - получаем заказ (основная логика с защитой от утечек)
func (u *Usecase) GetOrder(ctx context.Context, orderUID string) (*entity.Order, error) {
	//Ищем в кэше ID заказа
	data, exist := u.GetOrderFromCache(orderUID)
	if exist {
		//проверяем БД на случай если там нет заказа
		if err := u.CheckOrderFromCacheInDB(ctx, orderUID, exist); err != nil {
			log.Printf("Warning: failed to check order in DB: %v", err)
		}

		return data, nil
	}

	// Если нет в кэше ищем в БД
	data, err := u.GetOrderFromDB(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order from DB: %w", err)
	}

	// сохраняем в кэш
	if err := u.SaveOrderInCache(data); err != nil {
		log.Printf("Warning: failed to save order in cache: %v", err)
	}

	return data, nil
}

// HandleOrder - обработка заказа из Kafka с защитой от потери данных
func (u *Usecase) HandleOrder(ctx context.Context, order *entity.Order) error {
	// Проверяем что ID не пустой
	if order.OrderUID == "" {
		log.Printf("Invalid order: missing order_uid")
		return nil
	}

	//сохраняем в кэш
	if err := u.SaveOrderInCache(order); err != nil {
		return fmt.Errorf("failed to save order in cache: %w", err)
	}

	// сохраняем в БД
	return u.SaveOrderInDB(ctx, order)
}
