package grpc

import (
	"context"

	pb "github.com/Wuchinator/food-delivery/order-service/pkg/order_v1"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) CreateOrder(ctx context.Context,
	req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	return &pb.CreateOrderResponse{
		OrderId: 1231312,
		Status:  "Fake",
	}, nil
}
