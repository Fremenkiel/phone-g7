package admin

import (
	"log/slog"
	"net/http"

	"github.com/fremenkiel/phony-g7/internal/platform/logger"
)

type renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type Handler struct {
	renderer	renderer
}

func NewHandler(r renderer) *Handler {
	return &Handler{renderer: r}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if err := h.renderer.RenderPage(w, "index", nil); err != nil {
		logger.LogError(
			"table/handler",
			"error returned while rendering",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
