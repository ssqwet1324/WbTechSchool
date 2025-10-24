package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/wb-go/wbf/zlog"
)

// Config - структура конфига
type Config struct {
	DbName          string        `env:"DB_NAME"`
	DbUser          string        `env:"DB_USER"`
	DbPassword      string        `env:"DB_PASSWORD"`
	DbHost          string        `env:"DB_HOST"`
	DbPort          int           `env:"DB_PORT"`
	TimeZone        string        `env:"TIMEZONE"`
	RedisAddr       string        `env:"REDIS_ADDR"`
	RabbitURL       string        `env:"RABBIT_URL"`
	MaxRetries      int           `env:"MAX_RETRIES"`
	RetryDelay      time.Duration `env:"RETRY_DELAY"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
}

// New - конструктор
func New() (*Config, error) {
	var cfg Config
	_ = godotenv.Load(".env.example")

	cfg.DbName = os.Getenv("DB_NAME")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.DbPassword = os.Getenv("DB_PASSWORD")
	cfg.DbHost = os.Getenv("DB_HOST")

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		zlog.Logger.Error().Msg("error converting DB_PORT")
		return nil, fmt.Errorf("error converting DB_PORT: %w", err)
	}
	cfg.DbPort = dbPort

	cfg.TimeZone = os.Getenv("TIMEZONE")
	cfg.RedisAddr = os.Getenv("REDIS_ADDR")
	cfg.RabbitURL = os.Getenv("RABBIT_URL")

	maxRetriesInt, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		return nil, fmt.Errorf("error converting MAX_RETRIES: %w", err)
	}
	cfg.MaxRetries = maxRetriesInt

	retryDelayInt, err := strconv.Atoi(os.Getenv("RETRY_DELAY"))
	if err != nil {
		return nil, fmt.Errorf("error converting RETRY_DELAY: %w", err)
	}
	cfg.RetryDelay = time.Duration(retryDelayInt)

	maxOpenConnsInt, err := strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))
	if err != nil {
		return nil, fmt.Errorf("error converting MAX_OPEN_CONNS: %w", err)
	}
	cfg.MaxOpenConns = maxOpenConnsInt

	maxIdleConnsInt, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	if err != nil {
		return nil, fmt.Errorf("error converting MAX_IDLE_CONNS: %w", err)
	}
	cfg.MaxIdleConns = maxIdleConnsInt

	connMaxLifetimeInt, err := strconv.Atoi(os.Getenv("CONN_MAX_LIFETIME"))
	if err != nil {
		return nil, fmt.Errorf("error converting CONN_MAX_LIFETIME: %w", err)
	}
	cfg.ConnMaxLifetime = time.Duration(connMaxLifetimeInt)

	return &cfg, nil
}
