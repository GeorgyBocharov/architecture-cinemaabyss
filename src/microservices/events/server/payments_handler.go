package server

import (
	"fmt"
	"encoding/json"
	"net/http"

	"events/internal/entities"
)

type (
	PaymentsHandler struct {
		sender PaymentsSender
	}
	PaymentsSender interface {
		Send(payment entities.Payment) error
	}
)

func NewPaymentsHandler(sender PaymentsSender) *PaymentsHandler {
	return &PaymentsHandler{
		sender: sender,
	}
}

func (p *PaymentsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		p.createPayment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (p *PaymentsHandler) createPayment(w http.ResponseWriter, r *http.Request) {
	var payment entities.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("sending payment %v\n", payment)

	err := p.sender.Send(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	w.WriteHeader(http.StatusCreated)
}

