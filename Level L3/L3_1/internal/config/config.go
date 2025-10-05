package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config - структура конфига
type Config struct {
	DbName     string `env:"DB_NAME"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     int    `env:"DB_PORT"`
	TimeZone   string `env:"TIMEZONE"`
	RedisAddr  string `env:"REDIS_ADDR"`
	RabbitURL  string `env:"RABBIT_URL"`
	MaxRetries int    `env:"MAX_RETRIES"`
	RetryDelay int    `env:"RETRY_DELAY"`
}

// New - конструктор
func New() (*Config, error) {
	var cfg Config
	_ = godotenv.Load(".env")

	cfg.DbName = os.Getenv("DB_NAME")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.DbPassword = os.Getenv("DB_PASSWORD")
	cfg.DbHost = os.Getenv("DB_HOST")

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
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
	cfg.RetryDelay = retryDelayInt

	return &cfg, nil
}
