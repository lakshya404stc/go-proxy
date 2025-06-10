package server

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend interface {
	GetURL() *url.URL
	SetAlive(bool)
	IsAlive() bool
	GetProxy() *httputil.ReverseProxy
}

type ServerPool interface {
	AddBackend(Backend)
	NextPeer() Backend
	MarkBackendStatus(Backend, bool)
	GetServerCount() int
}

type backend struct {
	url   *url.URL
	proxy *httputil.ReverseProxy
	mu    sync.RWMutex
	alive bool
}

func NewBackend(url *url.URL, proxy *httputil.ReverseProxy) Backend {
	return &backend{
		url:   url,
		proxy: proxy,
		alive: true,
	}
}

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.alive = alive
	b.mu.Unlock()
}

func (b *backend) IsAlive() bool {
	b.mu.RLock()
	alive := b.alive
	b.mu.RUnlock()
	return alive
}

func (b *backend) GetProxy() *httputil.ReverseProxy {
	return b.proxy
}
