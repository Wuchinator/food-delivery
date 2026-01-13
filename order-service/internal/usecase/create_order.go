package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Wuchinator/food-delivery/order-service/internal/domain"
	"go.uber.org/zap"
)

type CreateOrderInput struct {
	UserID       int64
	RestaurantID int64
	Items        []CreateOrderItemInput
	Address      string
}

type CreateOrderItemInput struct {
	ProductID int64
	Quantity  int32
}

// Strcut of dependecies
type CreateOrderUseCase struct {
	repo   domain.OrderRepository
	logger *zap.Logger
	// ... kafka producer
}

func NewCreateOrderUseCase(repo domain.OrderRepository, logger *zap.Logger) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *CreateOrderUseCase) Exec(ctx context.Context, input CreateOrderInput) (int64, error) {

	if len(input.Items) == 0 {
		uc.logger.Error("goods can not be 0")
		return 0, errors.New("goods can not be 0")
	}

	orderItems := make([]domain.OrderItem, 0, len(input.Items))

	for _, item := range input.Items {
		fakePrice := rand.Intn(10) // Hardcode
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     int64(fakePrice),
		})
	}

	order, err := domain.NewOrder(input.UserID, input.RestaurantID, orderItems)
	if err != nil {
		uc.logger.Error("Failed to init new order", zap.Error(err))
	}

	orderID, err := uc.repo.Create(ctx, order)

	if err != nil {
		uc.logger.Error("Failed to create order", zap.Error(err))
		return 0, fmt.Errorf("Failed to create order %w", err)
	}

	return orderID, nil
}
