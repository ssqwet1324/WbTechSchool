package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/wb-go/wbf/zlog"
)

// Config - конфиг
type Config struct {
	DbName          string        `env:"DB_NAME" env-default:"postgres"`
	DbUser          string        `env:"DB_USER" env-default:"postgres"`
	DbPassword      string        `env:"DB_PASSWORD" env-default:"postgres"`
	DbHost          string        `env:"DB_HOST" env-default:"localhost"`
	DbPort          int           `env:"DB_PORT" env-default:"5432"`
	TimeZone        string        `env:"TIMEZONE" env-default:"UTC"`
	RedisAddr       string        `env:"REDIS_ADDR" env-default:"redis:6379"`
	MaxRetries      int           `env:"MAX_RETRIES" env-default:"3"`
	RetryDelay      time.Duration `env:"RETRY_DELAY" env-default:"5s"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" env-default:"30s"`
	BookRedisTTL    time.Duration `env:"BOOK_REDIS_TTL" env-default:"5m"`
}

// New - конструктор конфига
func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Error reading config")
		return nil, err
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
