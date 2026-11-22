package status


import (
	"net/http"
)

func NewRouter(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/status", handler.Status)
}
