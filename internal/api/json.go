package api

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) DecodeJSON(body io.ReadCloser, data any) error {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return err
	}
	defer func() {
		if err := body.Close(); err != nil {
			zap.L().Error("failed to close body", zap.Error(err))
		}
	}()

	return nil
}

func (r *Renderer) Json(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil || data == http.NoBody {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		zap.L().Error("failed to encode json", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *Renderer) Csv(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=report.csv")
	w.WriteHeader(status)
	if data == nil || data == http.NoBody {
		return
	}

	writer := csv.NewWriter(w)
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			zap.L().Error("failed to flush csv writer", zap.Error(err))
		}
	}()

	tableRows, ok := data.([][]string)
	if !ok {
		zap.L().Error("failed to cast data to [][]string")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, row := range tableRows {
		if err := writer.Write(row); err != nil {
			zap.L().Error("failed to write row", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (r *Renderer) JsonError(w http.ResponseWriter, err error) {
	var statusCode = http.StatusInternalServerError

	var apiError Error
	if ok := errors.As(err, &apiError); ok {
		statusCode = apiError.StatusCode()
	} else {
		zap.L().Error("unknown error", zap.Error(err))
		r.Json(w, statusCode, baseError{Msg: "unknown error"})
		return
	}

	r.Json(w, statusCode, apiError)
}
