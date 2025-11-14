package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/wb-go/wbf/zlog"
)

// ServiceConfig - конфиг
type ServiceConfig struct {
	DbName                string        `env:"DB_NAME" env-default:"postgres"`
	DbUser                string        `env:"DB_USER" env-default:"postgres"`
	DbPassword            string        `env:"DB_PASSWORD" env-default:"postgres"`
	DbHost                string        `env:"DB_HOST" env-default:"localhost"`
	DbPort                int           `env:"DB_PORT" env-default:"5432"`
	TimeZone              string        `env:"TIMEZONE" env-default:"UTC"`
	MaxRetries            int           `env:"MAX_RETRIES" env-default:"3"`
	RetryDelay            time.Duration `env:"RETRY_DELAY" env-default:"5s"`
	MaxOpenConns          int           `env:"MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns          int           `env:"MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime       time.Duration `env:"CONN_MAX_LIFETIME" env-default:"30m"`
	JWTSecret             string        `env:"JWT_SECRET" env-default:"supersecret"`
	PackageWithMigrations string        `env:"PACKAGE_WITH_MIGRATIONS" env-default:"./migrations"`
}

// New - конструктор конфига
func New() (*ServiceConfig, error) {
	var cfg ServiceConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Error reading config")
		return nil, err
	}

	zlog.Logger.Info().Any("cfg", cfg).Msg("config loaded successfully")

	return &cfg, nil
}
