package domain

import (
	"context"
	"time"
)

type KitchenStatus string

const (
	KitchenStatusAccepted  KitchenStatus = "ACCEPTED"
	KitchenStatusPreparing KitchenStatus = "PREPARING"
	KitchenStatusReady     KitchenStatus = "READY"
)

type KitchenOrder struct {
	ID           int64
	OrderID      int64
	RestaurantID int64
	Status       KitchenStatus
	Items        *KitchenItem
	CreatedAt    time.Time
}

type KitchenItem struct {
	ProductID int64
	Quantity  int32
}

type KitchenRepository interface {
	Create(ctx context.Context, order *KitchenOrder) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status KitchenStatus) error
}
