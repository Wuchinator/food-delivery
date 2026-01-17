package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	*pgxpool.Pool
	logger *zap.Logger
}

type Config struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	Timeout      time.Duration
}

func NewConn(cfg Config, logger *zap.Logger) (*DB, error) {
	db, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		logger.Error("Failed to connect DB", zap.Error(err))
		return nil, err
	}

	db.MaxConns = int32(cfg.MaxOpenConns)
	db.MinConns = int32(cfg.MaxIdleConns)
	db.MaxConnLifetime = cfg.Timeout

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, db)
	if err != nil {
		logger.Error("Failed to create pool conns", zap.Error(err))
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping databse", zap.Error(err))
		return nil, err
	}

	logger.Info("Pgx pool connection successfully started",
		zap.Int("max open conns", cfg.MaxOpenConns),
		zap.Int("max idle conns", cfg.MaxIdleConns))

	return &DB{
		Pool: pool,
	}, nil
}
