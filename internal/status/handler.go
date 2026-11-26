package status

import (
	"encoding/json"
	"net/http"

	"github.com/fremenkiel/phony-g7/internal/g7"
)

type Handler struct {
	emulator *g7.Emulator
}

func NewHandler(e *g7.Emulator) *Handler {
	return &Handler{emulator: e}
}

type response struct {
	Connected bool `json:"connected"`

	// Subscriptions are the event names that the host registered, for example "App.*".
	Subscriptions []string `json:"subscriptions"`

	// Screen is what the display of the device shows, for example the deny screen.
	Screen g7.Screen `json:"screen"`
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		Connected:     h.emulator.Connected.Load(),
		Subscriptions: h.emulator.Subscriptions(),
		Screen:        h.emulator.Screen(),
	})
}
