package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	segmentService service.ISegment
	r              *api.Renderer
}

func NewHandler(
	segmentService service.ISegment,
	renderer *api.Renderer,
) *Handler {
	return &Handler{
		segmentService: segmentService,
		r:              renderer,
	}
}

func (h *Handler) Mount(r *mux.Router) {
	segmentRouter := r.PathPrefix("/segment").Subrouter()
	segmentRouter.HandleFunc("", h.newSegment).Methods("POST")
	segmentRouter.HandleFunc("/{id}", h.deleteSegment).Methods("DELETE")
}
