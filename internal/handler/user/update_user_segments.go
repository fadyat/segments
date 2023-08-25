package user

import (
	"avito-internship-2023/internal/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) updateUserSegments(w http.ResponseWriter, r *http.Request) {
	var updateUserSegmentsDTO dto.UpdateUserSegments
	if err := h.r.DecodeJSON(r.Body, &updateUserSegmentsDTO); err != nil {
		h.r.JsonError(w, err)
		return
	}

	err := h.segmentService.UpdateUserSegments(r.Context(), mux.Vars(r)["id"], &updateUserSegmentsDTO)
	if err != nil {
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusOK, http.NoBody)
}
