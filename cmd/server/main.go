package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andro-kes/inventory_service/internal/inverr"
	"github.com/andro-kes/inventory_service/internal/logger"
	"github.com/andro-kes/inventory_service/internal/rpc"
	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg := logger.Config{
		Level:        "debug",
		Encoding:     "console",
		FileRotation: false,
		Development:  true,
	}
	if err := logger.Init(cfg); err != nil {
		panic("failed to init logger")
	}
	zl, err := logger.Logger()
	if err != nil {
		panic("wrong init logger")
	}
	defer zl.Sync()

	zl.Info("Start inventory service...")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB_URL is not found")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := NewPool(ctx, zl, dbURL)
	if err != nil {
		panic(err.Error())
	}
	defer pool.Close()

	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		panic("GRPC_ADDR must be set")
	}
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic("listen error: " + err.Error())
	}

	grpcServer := grpc.NewServer()
	inventoryService := rpc.NewInventoryService(ctx, pool)
	pb.RegisterInventoryServiceServer(grpcServer, inventoryService)

	serveErr := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			serveErr <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	
	select {
	case sig := <-shutdown:
		zl.Info("Server shutdown...", zap.Any("signal", sig))
	case err := <- serveErr:
		zl.Error(err.Error())
		panic("failed to start inventory service")
	}

	grpcServer.GracefulStop()
}

func NewPool(ctx context.Context, zl *zap.Logger, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		zl.Error(err.Error())
		return nil, inverr.InvalidPoolConfig
	}
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		zl.Error(err.Error())
		return nil, inverr.CreatePoolError
	}
	
	attempts := 5
	delay := time.Second
	for i := 0; i < attempts; i++ {
		if err := pool.Ping(ctx); err == nil {
			break
		}
		zl.Warn("failed to ping", zap.Any("delay", delay))
		time.Sleep(delay)
		delay *= 2
	}
	if err := pool.Ping(ctx); err != nil {
		zl.Error("failed to connect to pool")
		return nil, inverr.CreatePoolError
	}

	zl.Info("successfully connect to pool")
	return pool, nil
}
