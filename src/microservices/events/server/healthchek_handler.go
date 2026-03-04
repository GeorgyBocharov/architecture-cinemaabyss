package server

import (
	"encoding/json"
	"net/http"
)

type (
	HealthcheckHandler struct {
	}
)

func (p *HealthcheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}