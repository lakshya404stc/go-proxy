package handler

import (
	"net/http"

	"github.com/leonardo5621/golang-load-balancer/internal/server"
)

type contextKey string

const RETRY_ATTEMPTED = contextKey("retry-attempted")

type LoadBalancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

type loadBalancer struct {
	serverPool server.ServerPool
}

func NewLoadBalancer(serverPool server.ServerPool) LoadBalancer {
	return &loadBalancer{
		serverPool: serverPool,
	}
}

func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	peer := lb.serverPool.NextPeer()
	if peer == nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer.GetProxy().ServeHTTP(w, r)
}

func AllowRetry(r *http.Request) bool {
	if retryAttempted, ok := r.Context().Value(RETRY_ATTEMPTED).(bool); ok {
		return !retryAttempted
	}
	return true
}
