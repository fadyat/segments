package ping

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"net/http"
)

type Handler struct {
	r *api.Renderer
}

func NewHandler(
	renderer *api.Renderer,
) *Handler {
	return &Handler{
		r: renderer,
	}
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	h.r.Json(w, http.StatusOK, dto.Ping{
		Message: "pong",
	})
}
