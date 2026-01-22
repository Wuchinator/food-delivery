package grpc

import pb "github.com/Wuchinator/food-delivery/restaurant-service/pkg/restaurant_v1"

type Server struct {
	pb.UnimplementedRestaurantServiceServer
}
