package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wuchinator/food-delivery/order-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OrderRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewOrderRepository(pool *pgxpool.Pool, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		pool:   pool,
		logger: logger.Named("order_repository"),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryOrder := `
		INSERT INTO orders (user_id, restaurant_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var orderID int64
	err = tx.QueryRow(ctx,
		queryOrder,
		order.UserID, order.RestaurantID, order.Status, order.CreatedAt, order.UpdatedAt,
	).Scan(&orderID)

	if err != nil {
		r.logger.Error("failed to insert order", zap.Error(err))
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}

	queryItem := `
		INSERT INTO orders_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`

	batch := &pgx.Batch{}
	for _, item := range order.Items {
		batch.Queue(queryItem, orderID, item.ProductID, item.Quantity, item.Price)
	}

	results := tx.SendBatch(ctx, batch)

	for range order.Items {
		_, err := results.Exec()
		if err != nil {
			results.Close()
			r.logger.Error("failed to insert order items via batch", zap.Error(err))
			return 0, fmt.Errorf("failed to insert order items: %w", err)
		}
	}

	if err := results.Close(); err != nil {
		return 0, fmt.Errorf("failed to close batch results: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("successfully created order", zap.Int64("orderID", orderID))

	return orderID, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryOrder := `
		SELECT id, user_id, restaurant_id, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	var order domain.Order
	err = tx.QueryRow(ctx, queryOrder, id).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order with id %d not found", id)
		}
		r.logger.Error("failed to get order by id", zap.Int64("order_id", id), zap.Error(err))
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	queryItems := `
		SELECT product_id, quantity, price
		FROM orders_items
		WHERE order_id = $1
	`
	rows, err := tx.Query(ctx, queryItems, id)
	if err != nil {
		r.logger.Error("failed to get order items", zap.Int64("order_id", id), zap.Error(err))
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			r.logger.Error("failed to scan order item", zap.Error(err))
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error during iterating order items", zap.Error(err))
		return nil, fmt.Errorf("iterating order items: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit read-only transaction: %w", err)
	}

	return &order, nil
}
