package server

import (
	"fmt"
	"net/http"
)

type UrlBasedProxyHandlerProvider struct {
	handlerByPrefix map[string]ProxyHandler
	defaultHandler ProxyHandler
}

func NewUrlBasedProxyHandlerProvider(handlerByPrefix map[string]ProxyHandler, defaultHandler ProxyHandler) *UrlBasedProxyHandlerProvider {
	return &UrlBasedProxyHandlerProvider {
			handlerByPrefix : handlerByPrefix,
			defaultHandler : defaultHandler,
	}
}

func (p *UrlBasedProxyHandlerProvider) PorvideByRequest(r *http.Request) ProxyHandler {
	fmt.Printf("returning")
	handler, ok := p.handlerByPrefix[r.URL.Path]
	if ok {
		return handler
	}

	return p.defaultHandler
}