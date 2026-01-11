package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderGrpc "github.com/Wuchinator/food-delivery/order-service/internal/handler/grpc"
	pb "github.com/Wuchinator/food-delivery/order-service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	grpcServer := grpc.NewServer()
	orderHandler := orderGrpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051") // в будущем заменить на cfg.GRPCPort

	if err != nil {
		log.Fatal("Error init grpc listener", zap.Error(err))
	}

	go func() {
		// log.Info("Starting grpc server", zap.String("port", cfg.GRPCPort))
		log.Println("Starting grpc server")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("Error init grpc server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// log.Info("Shutting down grpc server")
	log.Println("Shutting down grpc server")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-stopped:
		// log.Info("GRPC server stopped")
		log.Println("GRPC server stopped")

	case <-ctx.Done():
		// log.Warn("shutdown grpc server timed out")
		grpcServer.Stop()
	}
	// log.Info("gRPC server stopped")
	log.Println("GRPC server stopped")

}
