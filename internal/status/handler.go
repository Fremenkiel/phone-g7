package status

import (
	"encoding/json"
	"net/http"

	"github.com/fremenkiel/phony-g7/internal/g7"
)

type Handler struct {
	emulator				*g7.Emulator
}

func NewHandler(e *g7.Emulator) *Handler {
	return &Handler{emulator: e}
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"connected": h.emulator.Connected.Load()})
}
