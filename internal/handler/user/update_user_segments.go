package user

import (
	"avito-internship-2023/internal/dto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) updateUserSegments(w http.ResponseWriter, r *http.Request) {
	var updateUserSegmentsDTO dto.UpdateUserSegments
	if err := h.r.DecodeJSON(r.Body, &updateUserSegmentsDTO); err != nil {
		h.r.JsonError(w, err)
		return
	}

	userID := mux.Vars(r)["id"]
	err := h.segmentService.UpdateUserSegments(r.Context(), userID, &updateUserSegmentsDTO)
	if err != nil {
		zap.L().Info("failed to update user segments", zap.Error(err))
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusNoContent, http.NoBody)
}
