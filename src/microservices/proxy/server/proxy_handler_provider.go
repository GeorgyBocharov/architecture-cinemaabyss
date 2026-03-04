package server

import (
	"fmt"
	"net/http"
	"strings"
)

type UrlBasedProxyHandlerProvider struct {
	handlerByPrefix map[string]HttpHandler
	defaultHandler  HttpHandler
}

func NewUrlBasedProxyHandlerProvider(handlerByPrefix map[string]HttpHandler, defaultHandler HttpHandler) *UrlBasedProxyHandlerProvider {
	return &UrlBasedProxyHandlerProvider{
		handlerByPrefix: handlerByPrefix,
		defaultHandler:  defaultHandler,
	}
}

func (p *UrlBasedProxyHandlerProvider) PorvideByRequest(r *http.Request) HttpHandler {
	handler, ok := p.handlerByPrefix[r.URL.Path]
	if ok {
		return handler
	}

	return p.defaultHandler
}

type CompositeHandlerProvider struct {
	isProxyEnabled bool
	healthURL      string
	proxyURLPrefix string

	proxyHandler HttpHandler
	healthHandle HttpHandler
	defaultHandler HttpHandler
}

type CompositeHandlerProviderOption func(*CompositeHandlerProvider) *CompositeHandlerProvider

func WithProxy(prefix string, proxyHandler HttpHandler) CompositeHandlerProviderOption {
	return func (c *CompositeHandlerProvider)  *CompositeHandlerProvider {
		c.isProxyEnabled = true
		c.proxyURLPrefix = prefix
		c.proxyHandler = proxyHandler

		return c
	}
}

func NewCompositeHandlerProvider(healthURL string, healthHandle, defaultHandler HttpHandler, options ...CompositeHandlerProviderOption) *CompositeHandlerProvider {
	r := &CompositeHandlerProvider{
		isProxyEnabled: false,
		healthURL:  healthURL,
		healthHandle:  healthHandle,
		defaultHandler:  defaultHandler,
	}
	for _, option := range options {
		r = option(r)
	}

	return r
}

func (p *CompositeHandlerProvider) PorvideByRequest(r *http.Request) HttpHandler {
	fmt.Printf("searching handler for url %s, proxyIsEnabled = %v\n", r.URL.Path, p.isProxyEnabled)
	if r.URL.Path == p.healthURL {
		fmt.Println("returning healthHandler")

		return p.healthHandle
	}
	if p.isProxyEnabled && strings.HasPrefix(r.URL.Path, p.proxyURLPrefix) {
		fmt.Println("returning proxyHandler")

		return p.proxyHandler
	}
	fmt.Println("returning defaultHandler")

	return p.defaultHandler
}
