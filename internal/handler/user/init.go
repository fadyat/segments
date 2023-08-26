package user

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
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{id}/segment", h.updateUserSegments).Methods("PUT")
	userRouter.HandleFunc("/{id}/segment", h.getUserSegments).Methods("GET")
	userRouter.HandleFunc("/{id}/segment/ttl", h.joinSegmentsWithTTL).Methods("POST")
}
