package action

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fremenkiel/phony-g7/internal/g7"
)

type service interface {
	CreateEvent(context.Context, string) (int64, error)
}

type Handler struct {
	emulator	*g7.Emulator
	service		service
}

func NewHandler(e *g7.Emulator, s service) *Handler {
	return &Handler{emulator: e, service: s}
}

func (h *Handler) Trigger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		http.Error(w, "invalid or oversized payload", http.StatusBadRequest)
		return
	}

	body = bytes.TrimSpace(body)
	if len(body) == 0 || body[0] != '{' || !json.Valid(body) {
		http.Error(w, "expected a JSON object", http.StatusBadRequest)
		return
	}

	lastId, err := h.service.CreateEvent(ctx, string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	select {
	case h.emulator.Wake <- struct{}{}:
	default:
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]int64{"id": lastId})
}
