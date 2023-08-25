package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			zap.L().Debug(
				"request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			)

			next.ServeHTTP(w, r)

			zap.L().Debug("response",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", time.Since(start)),
			)
		},
	)
}
