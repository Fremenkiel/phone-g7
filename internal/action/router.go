package action

import (
	"net/http"
)

func NewRouter(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /api/trigger", handler.Trigger)
}

