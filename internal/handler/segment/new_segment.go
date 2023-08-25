package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) newSegment(w http.ResponseWriter, r *http.Request) {
	var createSegmentDTO dto.CreateSegment
	if err := h.r.DecodeJSON(r.Body, &createSegmentDTO); err != nil {
		zap.L().Info("failed to decode request body", zap.Error(err))
		h.r.JsonError(w, api.NewBadRequestError(err.Error()))
		return
	}

	createdSegment, err := h.segmentService.NewSegment(r.Context(), &createSegmentDTO)
	if err != nil {
		zap.L().Info("failed to create segment", zap.Error(err))
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusCreated, createdSegment)
	zap.S().Infof("segment created: %s", createdSegment.ID)
}
