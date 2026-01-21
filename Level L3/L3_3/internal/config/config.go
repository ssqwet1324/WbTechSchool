package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - конфиг
type Config struct {
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
func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %s", err)
	}

	return &cfg, nil
}

// CreateDsn - создание адреса подключения к бд
func (cfg *Config) CreateDsn() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		strconv.Itoa(cfg.DbPort),
		cfg.DbName,
	)

	return dsn
}
