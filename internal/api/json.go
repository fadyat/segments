package api

import (
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
