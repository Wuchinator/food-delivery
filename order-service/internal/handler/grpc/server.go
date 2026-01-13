package grpc

import (
	"context"
	"fmt"

	"github.com/Wuchinator/food-delivery/order-service/internal/usecase"
	pb "github.com/Wuchinator/food-delivery/order-service/pkg/order_v1"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	usecase *usecase.CreateOrderUseCase
	logger  *zap.Logger
}

func NewServer(uc *usecase.CreateOrderUseCase, logger *zap.Logger) *Server {
	return &Server{
		usecase: uc,
		logger:  logger,
	}
}

func (s *Server) CreateOrder(ctx context.Context,
	req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	inputItems := make([]usecase.CreateOrderItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		inputItems = append(inputItems, usecase.CreateOrderItemInput{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	input := usecase.CreateOrderInput{
		UserID:       req.UserId,
		RestaurantID: req.RestaurantId,
		Items:        inputItems,
		Address:      req.DeliveryAddress,
	}

	id, err := s.usecase.Exec(ctx, input)
	if err != nil {
		s.logger.Error("Failed to exec order usecase", zap.Error(err))
		return nil, fmt.Errorf("Failed to exec order usecase %w", err) // REFACTOR ON CODES GRPC FOR VEST PRACTICE
	}

	s.logger.Info("Created order response", zap.Int64("Id", id))
	return &pb.CreateOrderResponse{
		OrderId: id,
		Status:  "Created", // hardcode.
	}, nil
}
