package server

import (
	"fmt"
	"encoding/json"
	"net/http"

	"events/internal/entities"
)

type (
	MoviesHandler struct {
		sender MovieSender
	}
	MovieSender interface {
		Send(movie entities.Movie) error
	}
)

func NewMoviesHandler(sender MovieSender) *MoviesHandler {
	return &MoviesHandler{
		sender: sender,
	}
}

func (p *MoviesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		p.createMovie(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (p *MoviesHandler) createMovie(w http.ResponseWriter, r *http.Request) {
	var movie entities.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("sending movie %v\n", movie)

	err := p.sender.Send(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	w.WriteHeader(http.StatusCreated)
	encodeSuccessStatus(w)
}

func encodeSuccessStatus(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
