package server

import (
	"bytes"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
)

type (
	ProxyHandler interface {
		Handle(w http.ResponseWriter, r *http.Request)
	}
	DualPercentageProxyHandler struct {
		firstURL   string
		secondURL  string
		percentage int
	}
	BasicProxyHandler struct {
		proxyURL string
	}
)

func NewDualPercentageProxyHandler(firstURL, secondURL string, percentage int) *DualPercentageProxyHandler {
	return &DualPercentageProxyHandler{
		firstURL:   firstURL,
		secondURL:  secondURL,
		percentage: percentage,
	}
}

func (p *DualPercentageProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	proxyRequest(p.resolveHost(), w, r)
}


func (p *DualPercentageProxyHandler) resolveHost() string {
	randValue := rand.IntN(100)
	if randValue < p.percentage {
		return p.firstURL
	}
	return p.secondURL
}

func NewBasicProxyHandler(proxyURL string) *BasicProxyHandler {
	return &BasicProxyHandler{
		proxyURL:   proxyURL,
	}
}

func (p *BasicProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	proxyRequest(p.proxyURL, w, r)
}

func proxyRequest(targetHost string, w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	targetURL := targetHost + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}
	log.Printf("proxing request to url %s\n", targetURL)

	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	copyHeaders(proxyReq.Header, r.Header)

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to proxy request: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read proxy response", http.StatusInternalServerError)
		return
	}

	copyHeaders(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
