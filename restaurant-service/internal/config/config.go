package config

// TODO: Разнести на файлы оставив лишь интерфейс для подключения

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	LoggerLevel string
	GRPCPort    string
	MetricsPort string
	Kafka       KafkaConfig
	Postgres    PostgresConfig
}

type PostgresConfig struct {
	User            string
	Password        string
	Host            string
	Port            string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifeTime time.Duration
	SSLmode         string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	TimeOut time.Duration
}

func Load() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "PRODUCTION" {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		LoggerLevel: getEnv("LOGGER_LEVEL", "DEBUG"),
		GRPCPort:    getEnv("GRPCPORT", "50051"),
		MetricsPort: getEnv("METRICS_PORT", "9091"),
	}

	cfg.Postgres = PostgresConfig{
		Host:            getEnv("POSTGRES_HOST", "postgres-main"),
		Port:            getEnv("POSTGRES_PORT", "5432"),
		Database:        getEnv("POSTGRES_DB", "delivery"),
		User:            getEnv("POSTGRES_USER", "user"),
		Password:        getEnv("POSTGRES_PASSWORD", "password"),
		MaxOpenConns:    getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 5),
		MaxConnLifeTime: getEnvAsDuration("POSTGRES_MAX_LIFE_TIME", 5*time.Minute),
		SSLmode:         getEnv("POSTGRES_SSL_MODE", "disable"),
	}

	brokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	cfg.Kafka = KafkaConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   getEnv("KAFKA_TOPIC", "restaurant"),
		GroupID: getEnv("GROUP_ID", "restaurant-group"),
		TimeOut: getEnvAsDuration("TIMEOUT", time.Second*30),
	}
	return cfg, nil
}

func getEnvAsDuration(Val string, defaultVal time.Duration) time.Duration {
	strVal := os.Getenv(Val)
	if val, err := time.ParseDuration(strVal); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsInt(Val string, defaultVal int) int {
	strValue := os.Getenv(Val)
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultVal
}

func getEnv(Val, defaultVal string) string {
	if value := os.Getenv(Val); value != "" {

		return value
	}

	return defaultVal
}
