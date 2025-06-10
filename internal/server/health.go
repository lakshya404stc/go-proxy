package server

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

func LaunchHealthCheck(ctx context.Context, pool ServerPool) {
	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go healthCheck(pool)
		}
	}
}

func healthCheck(pool ServerPool) {
	for i := 0; i < pool.GetServerCount(); i++ {
		backend := pool.NextPeer()
		if backend == nil {
			continue
		}

		alive := isBackendAlive(backend.GetURL())
		pool.MarkBackendStatus(backend, alive)
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(u.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
