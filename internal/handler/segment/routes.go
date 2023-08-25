package segment

import "github.com/gorilla/mux"

func (h *Handler) Mount(r *mux.Router) {
	segmentRouter := r.PathPrefix("/segment").Subrouter()
	segmentRouter.HandleFunc("", h.newSegment).Methods("POST")
}
