package segment

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

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
