package action

import (
	"net/http"
)

func NewRouter(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /api/trigger", handler.Trigger)
	mux.HandleFunc("POST /api/actions/card-read", handler.CardRead)
	mux.HandleFunc("POST /api/actions/barcode", handler.Barcode)
	mux.HandleFunc("POST /api/actions/pin", handler.Pin)
}
