package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	LogLevel    string
	GRPCPort    string
	Postgres    PostgresConfig
	Kafka       KafkaConfig
}
type PostgresConfig struct {
	Host            string
	Port            string
	Database        string
	User            string
	Password        string
	MaxOpenConns    int // as more as possible
	MaxIdleConns    int // as less as possible
	ConnMaxLifeTime time.Duration
	SSLMode         string
}

type KafkaConfig struct {
	Brokers         []string
	Topic           string
	ProducerTimeout time.Duration
	RequireAcks     int
}

func Load() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "production" {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		GRPCPort:    getEnv("ORDER_SERVICE_PORT", "50051"),
	}

	cfg.Postgres = PostgresConfig{
		Host:            getEnv("POSTGRES_HOST", "postgres-main"),
		Port:            getEnv("POSTGRES_PORT", "5432"),
		Database:        getEnv("POSTGRES_DB", "delivery"),
		User:            getEnv("POSTGRES_USER", "user"),
		Password:        getEnv("POSTGRES_PASSWORD", "password"),
		MaxOpenConns:    getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 5),
		ConnMaxLifeTime: getEnvAsDuration("POSTGRES_MAX_LIFE_TIME", 5*time.Minute),
		SSLMode:         getEnv("POSTGRES_SSL_MODE", "disable"),
	}

	brokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	cfg.Kafka = KafkaConfig{
		Brokers:         strings.Split(brokers, ","),
		Topic:           getEnv("KAFKA_TOPIC_ORDER", "user-order"),
		ProducerTimeout: getEnvAsDuration("KAFKA_PRODUCER_TIMEOUT", time.Second*15),
		RequireAcks:     getEnvAsInt("KAFKA_REQUIRED_ACKS", -1),
	}

	return cfg, nil
}

func (cfg *PostgresConfig) PostgresDSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
	return dsn
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}

	return defaultValue
}
