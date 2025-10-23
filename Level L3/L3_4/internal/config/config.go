package config

import (
	"os"
	"strconv"
	"time"

	"github.com/wb-go/wbf/zlog"
)

// ServiceConfig - конфиг
type ServiceConfig struct {
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
	KafkaGroupId        string        `env:"KAFKA_GROUP_ID"`
}

// New - конструктор конфига
func New() *ServiceConfig {
	s := &ServiceConfig{
		DbName:     getEnv("DB_NAME", "postgres"),
		DbUser:     getEnv("DB_USER", "postgres"),
		DbPassword: getEnv("DB_PASSWORD", "postgres"),
		DbHost:     getEnv("DB_HOST", "localhost"),
		TimeZone:   getEnv("TIMEZONE", "UTC"),

		//Kafka
		KafkaAddr:    getEnv("KAFKA_ADDR", "kafka:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "photos_topic"),
		KafkaGroupId: getEnv("KAFKA_GROUP_ID", "photo-consumer-group"),

		// 🔹 MinIO
		MinioEndpoint:       getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:      getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:      getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIoPublicEndpoint: getEnv("MINIO_PUBLIC_ENDPOINT", "http://localhost:9000"),
		BucketName:          getEnv("BUCKET_NAME", "photos"),
		MinioUseSSl:         getEnvBool("MINIO_USE_SSL", false),
	}

	// 🔹 Целые числа
	s.DbPort = getEnvInt("DB_PORT", 5432)
	s.MaxRetries = getEnvInt("MAX_RETRIES", 3)
	s.MaxOpenConns = getEnvInt("MAX_OPEN_CONNS", 10)
	s.MaxIdleConns = getEnvInt("MAX_IDLE_CONNS", 5)

	// 🔹 Время
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

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		b, err := strconv.ParseBool(val)
		if err == nil {
			return b
		}
	}
	return defaultVal
}
