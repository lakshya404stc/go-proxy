package server

import (
	"sync"
	"sync/atomic"
)

type roundRobinPool struct {
	backends []Backend
	current  uint64
	mu       sync.RWMutex
}

func NewServerPool(strategy string) (ServerPool, error) {
	switch strategy {
	case "round-robin":
		return &roundRobinPool{
			backends: make([]Backend, 0),
		}, nil
	default:
		return &roundRobinPool{
			backends: make([]Backend, 0),
		}, nil
	}
}

func (p *roundRobinPool) AddBackend(backend Backend) {
	p.mu.Lock()
	p.backends = append(p.backends, backend)
	p.mu.Unlock()
}

func (p *roundRobinPool) NextPeer() Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.backends) == 0 {
		return nil
	}

	next := atomic.AddUint64(&p.current, uint64(1))
	idx := next % uint64(len(p.backends))
	return p.backends[idx]
}

func (p *roundRobinPool) MarkBackendStatus(backend Backend, alive bool) {
	backend.SetAlive(alive)
}

func (p *roundRobinPool) GetServerCount() int {
	p.mu.RLock()
	count := len(p.backends)
	p.mu.RUnlock()
	return count
}
