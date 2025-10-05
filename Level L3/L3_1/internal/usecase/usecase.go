package usecase

import (
	"L3_1/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)

// RepositoryProvider интерфейс бд
type RepositoryProvider interface {
	CreateNotification(ctx context.Context, notify entity.Notify) error
	CheckStatusNotification(ctx context.Context, notifyID uuid.UUID) (string, error)
	DeleteNotification(ctx context.Context, notifyID uuid.UUID) error
	GetNotifications(ctx context.Context, userID string) ([]entity.Notify, error)
}

// CacheProvider - интерфейсы кеша
type CacheProvider interface {
	AddNotifyInCash(ctx context.Context, key string, notifyCash entity.NotifyCache) error
	GetNotifyInCash(ctx context.Context, key string) (entity.NotifyCache, error)
	DeleteNotifyInCash(ctx context.Context, key string) error
	GetAllNotifyKeys(ctx context.Context) ([]string, error)
}

// CacheKey - константа для ключа в кэше
const CacheKey = "notify:"

// UseCaseNotify - структура для бизнес-логики
type UseCaseNotify struct {
	repository RepositoryProvider
	cache      CacheProvider
}

// New - конструктор для usecase
func New(repository RepositoryProvider, cache CacheProvider) *UseCaseNotify {
	return &UseCaseNotify{
		repository: repository,
		cache:      cache,
	}
}

// generateNotifyID - создать id уведомления
func generateNotifyID() uuid.UUID {
	return uuid.New()
}

// BuildKey - сформировать ключ для редиса0
func BuildKey(notifyID string) string {
	return fmt.Sprintf("%s%s", CacheKey, notifyID)
}

// CreateNotification - сохранить уведомление
func (u *UseCaseNotify) CreateNotification(ctx context.Context, notify entity.Notify) (entity.NotifyCache, error) {
	notify.NotifyID = generateNotifyID()

	// добавляем уведомление в БД
	if err := u.repository.CreateNotification(ctx, notify); err != nil {
		return entity.NotifyCache{}, fmt.Errorf("CreateNotificiation: repository err %w", err)
	}

	// формируем объект для кэша
	notifyCache := entity.NotifyCache{
		NotifyID:  notify.NotifyID,
		Title:     notify.Title,
		Body:      notify.Body,
		EventTime: notify.SendingDate,
		Status:    notify.Status,
		UserID:    notify.UserID,
	}

	key := BuildKey(notify.NotifyID.String())

	// добавляем в кэш
	if err := u.cache.AddNotifyInCash(ctx, key, notifyCache); err != nil {
		zlog.Logger.Err(err).Msg("Failed to add notify in cache")
	}

	zlog.Logger.Info().Str("notifyID", notify.NotifyID.String()).Msg("Notification created")
	return notifyCache, nil
}

// CheckStatusNotification - проверяем статус уведомления
func (u *UseCaseNotify) CheckStatusNotification(ctx context.Context, notifyID string) (string, error) {
	data, err := u.cache.GetNotifyInCash(ctx, BuildKey(notifyID))
	if err == nil && data.Status {
		return "Notification sent", nil
	}

	notifyIDUUID, err := uuid.Parse(notifyID)
	if err != nil {
		return "Notification not sent", nil
	}

	// иначе проверим в БД
	status, err := u.repository.CheckStatusNotification(ctx, notifyIDUUID)
	if err != nil {
		return "", fmt.Errorf("CheckStatusNotification: %w", err)
	}
	return status, nil

}

// DeleteNotification - удаляем уведомление
func (u *UseCaseNotify) DeleteNotification(ctx context.Context, notifyID string) error {
	notifyIDUUID, err := uuid.Parse(notifyID)
	if err != nil {
		return fmt.Errorf("DeleteNotification: %w", err)
	}
	zlog.Logger.Info().Str("notifyID", notifyIDUUID.String()).Msg("Notification deleted")

	if err := u.repository.DeleteNotification(ctx, notifyIDUUID); err != nil {
		return fmt.Errorf("DeleteNotification: %w", err)
	}

	// удаляем из кэша (передаём сырой notifyID, ключ строится внутри)
	if err := u.DeleteNotifyInCash(ctx, notifyID); err != nil {
		zlog.Logger.Warn().Err(err).Msg("Failed to delete notify from cache")
	}

	zlog.Logger.Info().Str("notifyID", notifyID).Msg("Notification deleted")
	return nil
}

// GetNotifications - получить все уведомления пользователя
func (u *UseCaseNotify) GetNotifications(ctx context.Context, userID string) ([]entity.Notify, error) {
	return u.repository.GetNotifications(ctx, userID)
}

// GetNotifyInCash - получить уведомление из кэша
func (u *UseCaseNotify) GetNotifyInCash(ctx context.Context, notifyID string) (entity.NotifyCache, error) {
	key := BuildKey(notifyID)
	return u.cache.GetNotifyInCash(ctx, key)
}

// AddNotifyInCash - добавить уведомление в кэш
func (u *UseCaseNotify) AddNotifyInCash(ctx context.Context, notifyCache entity.NotifyCache) error {
	key := BuildKey(notifyCache.NotifyID.String())
	return u.cache.AddNotifyInCash(ctx, key, notifyCache)
}

// GetNearNotify - возвращает ближайшее по времени уведомление из кэша
func (u *UseCaseNotify) GetNearNotify(ctx context.Context) (*entity.NotifyCache, error) {
	// получаем все ключи
	keys, err := u.GetAllNotifyKeys(ctx)
	if err != nil {
		return nil, err
	}

	// проверяем есть ли ключи в кэше
	if len(keys) == 0 {
		return nil, errors.New("redis: nil")
	}

	var nearest *entity.NotifyCache
	for _, key := range keys {
		notify, err := u.cache.GetNotifyInCash(ctx, key)
		if err != nil {
			zlog.Logger.Err(err).Str("key", key).Msg("Ошибка при получении уведомления из кэша")
			continue
		}

		// ищем ближайшее событие
		if nearest == nil || notify.EventTime.Before(nearest.EventTime) {
			nearest = &notify
		}
	}

	if nearest == nil {
		return nil, errors.New("redis: nil")
	}

	// возвращаем ближайшее уведомление
	return nearest, nil
}

// GetAllNotifyKeys - получаем ключи из кэша
func (u *UseCaseNotify) GetAllNotifyKeys(ctx context.Context) ([]string, error) {
	return u.cache.GetAllNotifyKeys(ctx)
}

// DeleteNotifyInCash - удаляем уведомление из кеша
func (u *UseCaseNotify) DeleteNotifyInCash(ctx context.Context, notifyID string) error {
	key := BuildKey(notifyID)
	if err := u.cache.DeleteNotifyInCash(ctx, key); err != nil {
		return fmt.Errorf("DeleteNotifyInCash: %w", err)
	}

	return nil
}
