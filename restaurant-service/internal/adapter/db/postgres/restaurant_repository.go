package postgres

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

func (r *RestaurantRepository) CreateMenuItem(ctx context.Context, Menu *domain.MenuItem) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return 0, err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO menu (restaurant_id, product_id, name, price, description, is_available)
	 VALUES ($1, $2, $3, $4, $5, $6)
	 RETURNING id`

	var MenuID int64

	err = tx.QueryRow(ctx, query,
		Menu.RestaurantID, Menu.ProductID, Menu.Name,
		Menu.Price, Menu.Description, Menu.IsAvailable).Scan(&MenuID)

	if err != nil {
		r.logger.Error("Failed to insert data in menu", zap.Error(err))
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return 0, err
	}

	return MenuID, nil
}

func (r *RestaurantRepository) GetMenu(ctx context.Context, restaurantID int64) ([]domain.MenuItem, error) {
	return nil, nil
}

func (r *RestaurantRepository) UpdateMenu(ctx context.Context, item *domain.MenuItem) error {
	return nil
}

func (r *RestaurantRepository) DeleteMenu(ctx context.Context, restaurantID int64, itemID int64) error {
	return nil
}
