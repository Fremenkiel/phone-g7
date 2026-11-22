package event

import (
	"net/http"
)

func NewRouter(mux *http.ServeMux, handler *Handler, middlewares ...func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /api/events", handler.Events)
}
