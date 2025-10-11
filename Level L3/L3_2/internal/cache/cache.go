package cache

import (
	"context"
	"encoding/json"
	"errors"
	"shortener/internal/entity"
	"time"

	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

type Cache struct {
	redis redis.Client
}

func New(redis redis.Client) *Cache {
	return &Cache{redis: redis}
}

// AddShortUrlInCache - добавляем уведомление в кеш
func (c *Cache) AddShortUrlInCache(ctx context.Context, key string, notifyCash []byte, ttl time.Duration) error {
	return c.redis.SetWithExpiration(ctx, key, notifyCash, ttl)
}

// GetShortUrlFromCache - получить url из кеша
func (c *Cache) GetShortUrlFromCache(ctx context.Context, key string) (*entity.ShortenURL, error) {
	var urls entity.ShortenURL

	val, err := c.redis.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			zlog.Logger.Info().Msgf("no url in cache")
		}
		zlog.Logger.Warn().Err(err).Msgf("failed to get url from cache")
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &urls)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msgf("failed to unmarshal url from cache")
		return nil, err
	}

	return &urls, nil
}
