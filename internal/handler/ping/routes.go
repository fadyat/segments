package ping

import "github.com/gorilla/mux"

func (h *Handler) Mount(r *mux.Router) {
	r.HandleFunc("/ping", h.ping).Methods("GET")
}
