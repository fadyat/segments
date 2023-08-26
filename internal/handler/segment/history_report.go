package segment

import (
	"avito-internship-2023/internal/api"
	"net/http"
)

func (h *Handler) getHistoryReport(w http.ResponseWriter, r *http.Request) {
	timeRange := r.URL.Query().Get("time_range")
	responseFormat := api.ResponseFormat(r.URL.Query().Get("format"))
	if timeRange == "" {
		h.r.JsonError(w, api.NewBadRequestError("time_range is required"))
		return
	}

	if !responseFormat.IsValid() {
		h.r.JsonError(w, api.NewBadRequestError("invalid response format"))
		return
	}

	if responseFormat == api.ResponseFormatNone {
		responseFormat = api.ResponseFormatCSV
	}

	report, err := h.segmentService.GetHistoryReport(r.Context(), timeRange)
	if err != nil {
		h.r.JsonError(w, err)
		return
	}

	switch responseFormat {
	case api.ResponseFormatJSON:
		h.r.Json(w, http.StatusOK, report)
	case api.ResponseFormatCSV:
		h.r.Csv(w, http.StatusOK, report.ToRawTable())
	}
}
