package user

import (
	"avito-internship-2023/internal/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) joinSegmentsWithTTL(w http.ResponseWriter, r *http.Request) {
	var joinSegmentsWithTTLDTO []*dto.SegmentWithTTL
	if err := h.r.DecodeJSON(r.Body, &joinSegmentsWithTTLDTO); err != nil {
		h.r.JsonError(w, err)
		return
	}

	userID := mux.Vars(r)["id"]
	err := h.segmentService.JoinSegmentsWithTTL(r.Context(), userID, joinSegmentsWithTTLDTO)
	if err != nil {
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusNoContent, http.NoBody)
}
