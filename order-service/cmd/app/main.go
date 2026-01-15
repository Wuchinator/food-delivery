package main

import (
	"log"
	"time"

	"github.com/Wuchinator/food-delivery/order-service/internal/adapter/db/postgres"
	"github.com/Wuchinator/food-delivery/order-service/internal/adapter/kafka"
	"github.com/Wuchinator/food-delivery/order-service/internal/app"
	"github.com/Wuchinator/food-delivery/order-service/internal/app/database"
	"github.com/Wuchinator/food-delivery/order-service/internal/app/logger"
	"github.com/Wuchinator/food-delivery/order-service/internal/config"
	orderGrpc "github.com/Wuchinator/food-delivery/order-service/internal/handler/grpc"
	"github.com/Wuchinator/food-delivery/order-service/internal/usecase"
	pb "github.com/Wuchinator/food-delivery/order-service/pkg/order_v1"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

	kafka := kafka.NewProducer(kafka.Config{
		Brokers:         cfg.Kafka.Brokers,
		Topic:           cfg.Kafka.Topic,
		ProducerTimeout: cfg.Kafka.ProducerTimeout,
		RequireAcks:     cfg.Kafka.RequireAcks,
	}, log)

	defer kafka.Close()

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

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	)

	grpc_prometheus.Register(grpcServer)

	orderRepo := postgres.NewOrderRepository(db.Pool, log)
	uc := usecase.NewCreateOrderUseCase(orderRepo, log, kafka)
	orderHandler := orderGrpc.NewServer(uc, log)
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	App := app.NewApp(cfg, log, grpcServer)
	App.Run()
}
