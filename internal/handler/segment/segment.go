package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
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

func (h *Handler) deleteSegment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.segmentService.DeleteSegment(r.Context(), id)
	if err != nil {
		zap.L().Info("failed to delete segment", zap.Error(err))
		h.r.JsonError(w, err)
		return
	}

	h.r.Json(w, http.StatusNoContent, http.NoBody)
}
