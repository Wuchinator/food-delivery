package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	*pgxpool.Pool
	logger *zap.Logger
}

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime time.Duration
}

func NewConnection(config Config, logger *zap.Logger) (*DB, error) {
	db, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	db.MaxConns = int32(config.MaxOpenConns)
	db.MinConns = int32(config.MaxIdleConns)
	db.MaxConnLifetime = config.MaxConnLifetime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("faield to create pgxpool %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error("Failed to ping pool", zap.Error(err), zap.String("DSN", config.DSN))
		return nil, fmt.Errorf("Failed to ping pool %w", err)
	}

	logger.Info("Postgresql connected",
		zap.Int("max_open_conns", config.MaxOpenConns),
		zap.Int("max_idle_conns", config.MaxIdleConns))

	return &DB{
		Pool: pool,
	}, nil
}
