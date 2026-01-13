package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wuchinator/food-delivery/order-service/internal/adapter/db/postgres"
	"github.com/Wuchinator/food-delivery/order-service/internal/app/database"
	"github.com/Wuchinator/food-delivery/order-service/internal/app/logger"
	"github.com/Wuchinator/food-delivery/order-service/internal/config"
	orderGrpc "github.com/Wuchinator/food-delivery/order-service/internal/handler/grpc"
	"github.com/Wuchinator/food-delivery/order-service/internal/usecase"
	pb "github.com/Wuchinator/food-delivery/order-service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	log, err := logger.NewLogger(cfg.LogLevel, cfg.Environment)
	if err != nil {
		log.Fatal("Failed to init logger")
	}

	defer log.Sync()

	log = logger.WithService(log, "order-service")

	log.Info("Starting order service", zap.String("environment", cfg.Environment), zap.String("grpc port", cfg.GRPCPort))

	db, err := database.NewConnection(database.Config{
		DSN:             cfg.Postgres.PostgresDSN(),
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		MaxConnLifetime: cfg.Postgres.ConnMaxLifeTime,
	}, log)

	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	defer db.Close()

	orderRepo := postgres.NewOrderRepository(db.Pool, log)
	uc := usecase.NewCreateOrderUseCase(orderRepo, log)

	grpcServer := grpc.NewServer()
	orderHandler := orderGrpc.NewServer(uc, log)

	pb.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)

	if err != nil {
		log.Fatal("Error init grpc listener", zap.Error(err))
	}

	go func() {
		log.Info("Starting grpc server", zap.String("port", ":"+cfg.GRPCPort))
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("Error init grpc server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Shutting down grpc server")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-stopped:
		log.Info("GRPC server stopped")

	case <-ctx.Done():
		log.Warn("shutdown grpc server timed out")
		grpcServer.Stop()
	}
	log.Info("gRPC server stopped")

}
