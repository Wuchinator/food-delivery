package domain

import (
	"context"
	"time"
)

type OrderStatus string

const (
	OrderCreated   OrderStatus = "Created"
	OrperPaid      OrderStatus = "Paid"
	OrderCancelled OrderStatus = "Cancelled"
)

type Order struct {
	ID           int64
	UserID       int64
	RestaurantID int64
	Items        []OrderItem
	Status       OrderStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OrderItem struct {
	ProductID int64
	Quantity  int32
	Price     int64
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) (int64, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	// ...
}

func NewOrder(userID, restaurantID int64, items []OrderItem) (*Order, error) { // Пока без проверок

	now := time.Now()
	return &Order{
		UserID:       userID,
		RestaurantID: restaurantID,
		Items:        items,
		Status:       OrderCreated,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
