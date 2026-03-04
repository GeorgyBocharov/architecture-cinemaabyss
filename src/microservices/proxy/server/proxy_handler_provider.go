package server

import (
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
	return &CompositeHandlerProvider{
		isProxyEnabled: false,
		healthURL:  healthURL,
		healthHandle:  healthHandle,
		defaultHandler:  defaultHandler,
	}
}

func (p *CompositeHandlerProvider) PorvideByRequest(r *http.Request) HttpHandler {
	if r.URL.Path == p.healthURL {
		return p.healthHandle
	}
	if p.isProxyEnabled && strings.HasPrefix(r.URL.Path, p.proxyURLPrefix) {
		return p.proxyHandler
	}

	return p.defaultHandler
}
