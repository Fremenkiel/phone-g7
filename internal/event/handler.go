package event

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type service interface {
	GetEvents(context.Context) ([]Event, error) 
}

type Handler struct {
	service		service
}

func NewHandler(s service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}
	events, err := h.service.GetEvents(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("event response: %v", err)
	}
}
