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
	emulator *g7.Emulator
	service  service
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

	h.wake()

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]int64{"id": lastId})
}

// CardRead imitates a data carrier that a person presents to the reader.
func (h *Handler) CardRead(w http.ResponseWriter, r *http.Request) {
	h.queueIdent(w, r, g7.NewCardReadEvent)
}

// Barcode imitates a barcode that a person presents to the barcode reader.
func (h *Handler) Barcode(w http.ResponseWriter, r *http.Request) {
	h.queueIdent(w, r, g7.NewBarcodeEvent)
}

// Pin imitates a PIN that a person enters at the device.
func (h *Handler) Pin(w http.ResponseWriter, r *http.Request) {
	h.queueIdent(w, r, g7.NewPinEvent)
}

// queueIdent writes the event App.CardIdent to the database with the status "queued" and wakes
// the emulator. The emulator sends the event when the host is registered to it.
func (h *Handler) queueIdent(w http.ResponseWriter, r *http.Request, build func(id string) (g7.Command, error)) {
	var request IdentRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&request); err != nil {
		http.Error(w, `expected a JSON object, for example {"id":"12345"}`, http.StatusBadRequest)
		return
	}

	command, err := build(request.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastId, err := h.service.CreateEvent(r.Context(), string(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wake()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(IdentResponse{
		ID:         lastId,
		Subscribed: h.emulator.SubscribedTo(command.CommandType),
	})
}

func (h *Handler) wake() {
	select {
	case h.emulator.Wake <- struct{}{}:
	default:
	}
}
