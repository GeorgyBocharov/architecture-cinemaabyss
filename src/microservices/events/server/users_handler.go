package server

import (
	"fmt"
	"encoding/json"
	"net/http"

	"events/internal/entities"
)

type (
	UsersHandler struct {
		sender UserSender
	}
	UserSender interface {
		Send(user entities.User) error
	}
)

func NewUsersHandler(sender UserSender) *UsersHandler {
	return &UsersHandler{
		sender: sender,
	}
}

func (p *UsersHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		p.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (p *UsersHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("sending user %v\n", user)

	err := p.sender.Send(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	w.WriteHeader(http.StatusCreated)
	encodeSuccessStatus(w)
}

