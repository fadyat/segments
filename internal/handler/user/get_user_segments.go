package user

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) getUserSegments(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	segments, err := h.segmentService.GetActiveUserSegments(r.Context(), userID)
	if err != nil {
		zap.L().Info("failed to get user segments", zap.Error(err))
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusOK, segments)
}
