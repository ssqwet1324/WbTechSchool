package usecase

import (
	"context"
	"errors"
	"sales_tracker/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)

type RepositoryProvider interface {
	AddItems(ctx context.Context, item entity.Item) error
	GetItems(ctx context.Context, items entity.GetItems) ([]entity.Item, error)
	UpdateItems(ctx context.Context, item entity.Item) error
	DeleteItems(ctx context.Context, itemID string) error
	GetAnalytics(ctx context.Context, analytics entity.GetItemsFromAnalytics) (*entity.AnalyticsResult, error)
}

type UseCase struct {
	repo RepositoryProvider
}

func New(repo RepositoryProvider) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func generateNewID() string {
	return uuid.New().String()
}

// AddItems - добавить запись
func (uc *UseCase) AddItems(ctx context.Context, item entity.NewItem) (string, error) {
	if item.Amount < 0 {
		zlog.Logger.Error().Msg("Introduced with negative amount")
		return "", errors.New("amount must be greater than zero")
	}

	newID := generateNewID()

	now := time.Now()

	itemStruct := entity.Item{
		ID:        newID,
		Title:     item.Title,
		Amount:    item.Amount,
		Date:      item.Date,
		Category:  item.Category,
		CreatedAt: now,
	}

	err := uc.repo.AddItems(ctx, itemStruct)
	if err != nil {
		return "", err
	}

	return newID, nil
}

// GetItems - получаем записи
func (uc *UseCase) GetItems(ctx context.Context, items entity.GetItems) ([]entity.Item, error) {
	return uc.repo.GetItems(ctx, items)
}

// UpdateItems - обновляем данные в записи
func (uc *UseCase) UpdateItems(ctx context.Context, item entity.NewItem, itemID string) error {
	if item.Amount < 0 {
		zlog.Logger.Error().Msg("Introduced with negative amount")
		return errors.New("amount must be greater than zero")
	}

	itemStruct := entity.Item{
		ID:        itemID,
		Title:     item.Title,
		Amount:    item.Amount,
		Date:      item.Date,
		Category:  item.Category,
		UpdatedAt: time.Now(),
	}

	err := uc.repo.UpdateItems(ctx, itemStruct)
	if err != nil {
		return err
	}

	return nil
}

// DeleteItems - удалить запись
func (uc *UseCase) DeleteItems(ctx context.Context, itemID string) error {
	err := uc.repo.DeleteItems(ctx, itemID)
	if err != nil {
		return err
	}

	return nil
}

// GetAnalytics - получаем аналитику по записям
func (uc *UseCase) GetAnalytics(ctx context.Context, analytics entity.GetItemsFromAnalytics) (*entity.AnalyticsResult, error) {
	data, err := uc.repo.GetAnalytics(ctx, analytics)
	if err != nil {
		return nil, err
	}

	return data, nil
}
