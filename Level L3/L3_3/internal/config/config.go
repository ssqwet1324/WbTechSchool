package config

import (
	"os"
	"strconv"
	"time"

	"github.com/wb-go/wbf/zlog"
)

// ServiceConfig - конфиг
type ServiceConfig struct {
	DbName          string        `env:"DB_NAME"`
	DbUser          string        `env:"DB_USER"`
	DbPassword      string        `env:"DB_PASSWORD"`
	DbHost          string        `env:"DB_HOST"`
	DbPort          int           `env:"DB_PORT"`
	TimeZone        string        `env:"TIMEZONE"`
	MaxRetries      int           `env:"MAX_RETRIES"`
	RetryDelay      time.Duration `env:"RETRY_DELAY"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
}

// New - конструктор конфига
func New() *ServiceConfig {
	s := &ServiceConfig{
		DbName:     getEnv("DB_NAME", "postgres"),
		DbUser:     getEnv("DB_USER", "postgres"),
		DbPassword: getEnv("DB_PASSWORD", "postgres"),
		DbHost:     getEnv("DB_HOST", "localhost"),
		TimeZone:   getEnv("TIMEZONE", "UTC"),
	}

	// Целые числа
	s.DbPort = getEnvInt("DB_PORT", 5432)
	s.MaxRetries = getEnvInt("MAX_RETRIES", 3)
	s.MaxOpenConns = getEnvInt("MAX_OPEN_CONNS", 10)
	s.MaxIdleConns = getEnvInt("MAX_IDLE_CONNS", 5)

	// Время
	s.RetryDelay = getEnvDuration("RETRY_DELAY", 5*time.Second)
	s.ConnMaxLifetime = getEnvDuration("CONN_MAX_LIFETIME", 30*time.Second)

	zlog.Logger.Info().Msg("config loaded successfully")

	return s
}

// Вспомогательные функции
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

// getEnvInt - получить int значение
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}

	return defaultVal
}

// getEnvDuration получить временное значение
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}

	return defaultVal
}
