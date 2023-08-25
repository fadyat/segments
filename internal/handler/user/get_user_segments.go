package user

import (
	"avito-internship-2023/internal/dto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) getUserSegments(w http.ResponseWriter, r *http.Request) {
	status := dto.NewUserSegmentStatus(r.URL.Query().Get("status"))
	userID := mux.Vars(r)["id"]

	segments, err := h.segmentService.GetUserSegments(r.Context(), status, userID)
	if err != nil {
		zap.L().Info("failed to get user segments", zap.Error(err))
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusOK, segments)
}
