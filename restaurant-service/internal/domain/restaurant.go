package domain

import "context"

type Restaurant struct {
	ID      int64
	Name    string
	Address string
}

type MenuItem struct {
	ID           int64
	ProductID    int64
	RestaurantID int64
	Price        int64
	Name         string
	Description  string
	IsAvailable  bool
}

type RestaurantRepository interface {
	GetMenu(ctx context.Context, restaurantID int64) ([]MenuItem, error)
	UpdateMenu(ctx context.Context, item *MenuItem) error
	CreateMenuItem(ctx context.Context, Menu *MenuItem) (int64, error)
	DeleteMenu(ctx context.Context, restaurantID int64, itemID int64) error
}
