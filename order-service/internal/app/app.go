package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wuchinator/food-delivery/order-service/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	grpcServer *grpc.Server
	httpServer *http.Server
}

func NewApp(cfg *config.Config,
	logger *zap.Logger,
	grpcServer *grpc.Server) *App {

	httpServer := &http.Server{
		Addr:              ":" + cfg.PrometheusPort,
		Handler:           promhttp.Handler(),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return &App{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
		httpServer: httpServer,
	}
}

func (a *App) Run() {
	go func() {
		a.logger.Info("Starting metrics server", zap.String("addr", a.httpServer.Addr))
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.logger.Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	listener, err := net.Listen("tcp", ":"+a.cfg.GRPCPort)
	if err != nil {
		a.logger.Fatal("Failed to create grpc listener", zap.Error(err))
	}

	go func() {
		a.logger.Info("Starting grpc server", zap.String("addr", a.cfg.GRPCPort))
		if err := a.grpcServer.Serve(listener); err != nil {
			a.logger.Fatal("Failed to start grpc server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down app...")

	a.Stop()

}

func (a *App) Stop() {

	const timeOut = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	a.logger.Info("Stopping grpc server...")
	a.grpcServer.GracefulStop()

	a.logger.Info("Stoppong HTTP server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Warn("HTTP server shutdown error", zap.Error(err))
	}

	a.logger.Info("Application stopped")
}
