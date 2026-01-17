package domain

import (
	"context"
	"time"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurantOrders *RestaurantOrders) (int64, error)
	Read(ctx context.Context, ID int64) (bool, error)
}
type orderStatus string

const (
	Accepted  orderStatus = "Accepted"
	Ready     orderStatus = "Ready"
	Cancelled orderStatus = "Cancelled"
)

type RestaurantOrders struct {
	ID            int64
	OrderID       int64
	RestaurantIdD int64
	Status        orderStatus
	Items         []RestaurantOrderItems
	CreatedAt     time.Time
	AcceptedAt    time.Time
	ReadyAt       time.Time
	CancelledAt   time.Time
	UpdatedAt     time.Time
}

type RestaurantOrderItems struct {
	ProductID int64
	Quantity  int32
	Price     int64
}

func NewRestaurantOrders(OrderID, RestaurantID int64, status orderStatus, items []RestaurantOrderItems) *RestaurantOrders {
	now := time.Now()
	return &RestaurantOrders{
		OrderID:       OrderID,
		RestaurantIdD: RestaurantID,
		Status:        status,
		Items:         items,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
