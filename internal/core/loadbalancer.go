package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/leonardo5621/golang-load-balancer/internal/config"
	"github.com/leonardo5621/golang-load-balancer/internal/handler"
	"github.com/leonardo5621/golang-load-balancer/internal/server"
)

type LoadBalancer struct {
	config       *config.Config
	logger       *zap.Logger
	serverPool   server.ServerPool
	loadBalancer handler.LoadBalancer
}

func NewLoadBalancer(config *config.Config, logger *zap.Logger) (*LoadBalancer, error) {
	serverPool, err := server.NewServerPool(config.GetStrategy())
	if err != nil {
		return nil, err
	}

	return &LoadBalancer{
		config:       config,
		logger:       logger,
		serverPool:   serverPool,
		loadBalancer: handler.NewLoadBalancer(serverPool),
	}, nil
}

func (lb *LoadBalancer) setupBackends() error {
	for _, u := range lb.config.Backends {
		endpoint, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("failed to parse URL %s: %w", u, err)
		}

		rp := httputil.NewSingleHostReverseProxy(endpoint)
		backendServer := server.NewBackend(endpoint, rp)

		rp.ErrorHandler = lb.createErrorHandler(backendServer)
		lb.serverPool.AddBackend(backendServer)
	}
	return nil
}

func (lb *LoadBalancer) createErrorHandler(backendServer server.Backend) func(http.ResponseWriter, *http.Request, error) {
	return func(writer http.ResponseWriter, request *http.Request, e error) {
		lb.logger.Error("error handling the request",
			zap.String("host", backendServer.GetURL().Host),
			zap.Error(e),
		)
		backendServer.SetAlive(false)

		if !handler.AllowRetry(request) {
			lb.logger.Info(
				"Max retry attempts reached, terminating",
				zap.String("address", request.RemoteAddr),
				zap.String("path", request.URL.Path),
			)
			http.Error(writer, "Service not available", http.StatusServiceUnavailable)
			return
		}

		lb.logger.Info(
			"Attempting retry",
			zap.String("address", request.RemoteAddr),
			zap.String("URL", request.URL.Path),
			zap.Bool("retry", true),
		)
		lb.loadBalancer.Serve(
			writer,
			request.WithContext(
				context.WithValue(request.Context(), handler.RETRY_ATTEMPTED, true),
			),
		)
	}
}

func (lb *LoadBalancer) Start(ctx context.Context) error {
	if err := lb.setupBackends(); err != nil {
		return err
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", lb.config.Port),
		Handler: http.HandlerFunc(lb.loadBalancer.Serve),
	}

	go server.LaunchHealthCheck(ctx, lb.serverPool)

	go func() {
		<-ctx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	lb.logger.Info(
		"Load Balancer started",
		zap.Int("port", lb.config.Port),
	)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
