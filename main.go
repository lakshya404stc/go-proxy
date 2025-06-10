package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/leonardo5621/golang-load-balancer/internal/config"
	"github.com/leonardo5621/golang-load-balancer/internal/core"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	lb, err := core.NewLoadBalancer(cfg, logger)
	if err != nil {
		logger.Fatal("failed to create load balancer", zap.Error(err))
	}

	if err := lb.Start(ctx); err != nil {
		logger.Fatal("failed to start load balancer", zap.Error(err))
	}
}
