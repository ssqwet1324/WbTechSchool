package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - конфиг
type Config struct {
	DbName              string        `env:"DB_NAME"`
	DbUser              string        `env:"DB_USER"`
	DbPassword          string        `env:"DB_PASSWORD"`
	DbHost              string        `env:"DB_HOST"`
	DbPort              int           `env:"DB_PORT"`
	TimeZone            string        `env:"TIMEZONE"`
	MaxRetries          int           `env:"MAX_RETRIES"`
	RetryDelay          time.Duration `env:"RETRY_DELAY"`
	ConnMaxLifetime     time.Duration `env:"CONN_MAX_LIFETIME"`
	MinioEndpoint       string        `env:"MINIO_ENDPOINT"`
	MinioAccessKey      string        `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey      string        `env:"MINIO_SECRET_KEY"`
	MinioUseSSl         bool          `env:"MINIO_USE_SSL"`
	MinIoPublicEndpoint string        `env:"MINIO_PUBLIC_ENDPOINT"`
	BucketName          string        `env:"BUCKET_NAME"`
	MaxOpenConns        int           `env:"MAX_OPEN_CONNS"`
	MaxIdleConns        int           `env:"MAX_IDLE_CONNS"`
	KafkaAddr           string        `env:"KAFKA_ADDR"`
	KafkaTopic          string        `env:"KAFKA_TOPIC"`
	KafkaGroupID        string        `env:"KAFKA_GROUP_ID"`
	TimeOfLiveURL       time.Duration `env:"TIME_OF_LIVE_URL"`
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
