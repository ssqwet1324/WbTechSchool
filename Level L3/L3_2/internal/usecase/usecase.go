package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"shortener/internal/entity"
	"strings"
	"time"

	"github.com/wb-go/wbf/zlog"
)

// RepositoryProvider - интерфейс репы
type RepositoryProvider interface {
	AddShortURL(ctx context.Context, urls entity.ShortenURL) error
	GetShortURL(ctx context.Context, url string) (string, error)
	AddAnalytics(ctx context.Context, urls entity.ShortenURLAnalytics) error
	GetAnalytics(ctx context.Context, shortURL string) (*entity.ShortenURLAnalytics, error)
	ExistsShortURL(ctx context.Context, shortURL string) (bool, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
}

// CacheProvider - интерфейс кеша
type CacheProvider interface {
	AddShortURLInCache(ctx context.Context, key string, notifyCash []byte, ttl time.Duration) error
	GetShortURLFromCache(ctx context.Context, key string) (*entity.ShortenURL, error)
}

const (
	// cacheKey - ключ редиса
	cacheKey = "short_url:"

	// countPopular - количество переход за день для популярной ссылки
	countPopular = 3

	// TTLPopularURL - время жизни популярной ссылки в кеше
	TTLPopularURL = 7 * 24 * time.Hour
)

// UseCase - бизнес логика
type UseCase struct {
	repository RepositoryProvider
	cache      CacheProvider
}

// New - конструктор для usecase
func New(repository RepositoryProvider, cache CacheProvider) *UseCase {
	return &UseCase{
		repository: repository,
		cache:      cache,
	}
}

// buildRedisKey - билдим ключ для редиса
func buildRedisKey(key string) string {
	return fmt.Sprintf("%s%s", cacheKey, key)
}

// generateShortenURL - создаем короткую ссылку
func generateShortenURL(longURL string) string {
	baseChars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	hash := sha256.Sum256([]byte(longURL))

	hashInt := new(big.Int).SetBytes(hash[:])

	var result strings.Builder
	base := big.NewInt(62)

	for hashInt.Sign() > 0 {
		remainder := new(big.Int)
		hashInt.DivMod(hashInt, base, remainder)
		result.WriteByte(baseChars[remainder.Int64()])
	}

	return result.String()[:6]
}

// AddShortURL - добавляем сокращенную ссылку в бд
func (uc *UseCase) AddShortURL(ctx context.Context, urls entity.ShortenURL) (string, error) {
	// создаем короткую ссылку
	shortURL := generateShortenURL(urls.OriginalURL)
	urls.ShortURL = shortURL

	// проверяем редис на наличие такой
	if cachedURL, err := uc.GetShortURLFromCache(ctx, buildRedisKey(shortURL)); err == nil && cachedURL != "" {
		zlog.Logger.Info().Msgf("Shorten URL already exists: %s", shortURL)
		return cachedURL, nil
	}

	// тут проверяем есть ли такой в бд
	exist, err := uc.repository.ExistsShortURL(ctx, urls.ShortURL)
	if err != nil {
		return "", err
	}
	if exist {
		return shortURL, nil
	}

	// сохраняем в бд
	if err := uc.repository.AddShortURL(ctx, urls); err != nil {
		zlog.Logger.Warn().Err(err).Str("url", urls.OriginalURL).Msg("UseCase: Failed to add url")
		return "", err
	}

	return shortURL, nil
}

// GetOriginalURL - получить оригинальный Url
func (uc *UseCase) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	return uc.repository.GetOriginalURL(ctx, shortURL)
}

// AddAnalytics - делаем аналитику + если ссылка популярная кидаем ее в кеш
func (uc *UseCase) AddAnalytics(ctx context.Context, urls entity.ShortenURLAnalytics) error {
	// Добавляем аналитику в базу
	if err := uc.repository.AddAnalytics(ctx, urls); err != nil {
		zlog.Logger.Warn().Err(err).Msg("AddAnalytics: Failed to add analytics")
		return err
	}

	// Получаем актуальные данные
	data, err := uc.GetAnalytics(ctx, urls.ShortURL)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("AddAnalytics: Failed to get analytics")
		return err
	}

	// проверяем есть ли популярный в кеше
	if cachedURL, err := uc.GetShortURLFromCache(ctx, buildRedisKey(urls.ShortURL)); err == nil && cachedURL != "" {
		return nil
	}

	popular, err := uc.ParseAnalytics(ctx, *data)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("AddAnalytics: Failed to parse analytics")
		return err
	}

	if popular {
		zlog.Logger.Info().Str("url", urls.ShortURL).Msg("AddAnalytics: Popular URL added in cache")
	}

	return nil
}

// GetAnalytics - получаем аналитику
func (uc *UseCase) GetAnalytics(ctx context.Context, shortURL string) (*entity.ShortenURLAnalytics, error) {
	return uc.repository.GetAnalytics(ctx, shortURL)
}

// ParseAnalytics - проверяем популярность URL за день
func (uc *UseCase) ParseAnalytics(ctx context.Context, analytics entity.ShortenURLAnalytics) (bool, error) {
	var shorten entity.ShortenURL
	now := time.Now()
	date := now.Format("2006-01-02")

	// проверяем ссылку на популярность за день
	if val, exist := analytics.ByDay[date]; exist && val > countPopular {
		shorten.ShortURL = analytics.ShortURL
		if err := uc.AddShortURLInCache(ctx, buildRedisKey(analytics.ShortURL), shorten); err != nil {
			zlog.Logger.Error().Err(err).Str("url", analytics.ShortURL).Msg("UseCase: Failed to add url in cache")
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// AddShortURLInCache - добавляем данные в кеш
func (uc *UseCase) AddShortURLInCache(ctx context.Context, key string, notifyCash entity.ShortenURL) error {
	data, err := json.Marshal(notifyCash)
	if err != nil {
		return err
	}

	return uc.cache.AddShortURLInCache(ctx, key, data, TTLPopularURL)
}

// GetShortURLFromCache - получить данные с кеша
func (uc *UseCase) GetShortURLFromCache(ctx context.Context, key string) (string, error) {
	data, err := uc.cache.GetShortURLFromCache(ctx, key)
	if err != nil {
		zlog.Logger.Info().Err(err).Str("key", key).Msg("GetShortURLFromCache")
		return "", err
	}

	return data.ShortURL, nil
}
