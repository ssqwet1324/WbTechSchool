package cache

import (
	"L3_1/internal/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

// Cache - структура кеша
type Cache struct {
	cache *redis.Client
}

// New конструктор для кеша
func New(redisClient *redis.Client) *Cache {
	return &Cache{
		cache: redisClient,
	}
}

// AddNotifyInCash - добавляем уведомление в кеш
func (c *Cache) AddNotifyInCash(ctx context.Context, key string, notifyCash entity.NotifyCache) error {
	data, err := json.Marshal(notifyCash)
	if err != nil {
		return err
	}

	return c.cache.Set(ctx, key, data)
}

// GetNotifyInCash - получить уведомление из кеша
func (c *Cache) GetNotifyInCash(ctx context.Context, key string) (entity.NotifyCache, error) {
	var notifyCash entity.NotifyCache

	val, err := c.cache.Get(ctx, key)
	if errors.Is(err, redis.NoMatches) {
		zlog.Logger.Info().Msgf("No notify in cash")
		return notifyCash, nil
	} else if err != nil {
		return notifyCash, fmt.Errorf("get notify in cash err: %w", err)
	}

	err = json.Unmarshal([]byte(val), &notifyCash)
	if err != nil {
		return entity.NotifyCache{}, fmt.Errorf("unmarshal notify in cash err: %w", err)
	}

	return notifyCash, nil
}

// DeleteNotifyInCash удаляет уведомление из кэша по ключу.
func (c *Cache) DeleteNotifyInCash(ctx context.Context, key string) error {
	if err := c.cache.Del(ctx, key); err != nil {
		return fmt.Errorf("delete notify in cash err: %v", err)
	}

	return nil
}

// GetAllNotifyKeys - получаем ключи уведомлений
func (c *Cache) GetAllNotifyKeys(ctx context.Context) ([]string, error) {
	var (
		cursor uint64 = 0
		keys   []string
	)

	for {
		scanKeys, newCursor, err := c.cache.Scan(ctx, cursor, "notify:*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("GetAllNotifyKeys: redis scan error: %w", err)
		}

		keys = append(keys, scanKeys...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
