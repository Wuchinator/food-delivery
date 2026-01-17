package database

import (
	"context"
	"restaurant/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type RestaurantRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewRestaurantRepository(pool *pgxpool.Pool, logger *zap.Logger) *RestaurantRepository {
	return &RestaurantRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *RestaurantRepository) Create(ctx context.Context, restaurantOrders *domain.RestaurantOrders) (int64, error) {
	_, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("Failed to start transaction", zap.Error(err))
		return 0, err
	}

	//query := `INSERT INTO restaurant_orders VALUES ($1 $2 $3)`
	return 0, nil
}
func (r *RestaurantRepository) Read(ctx context.Context, ID int64) (bool, error) {
	return false, nil
}
