package api

import (
	"avito-internship-2023/internal/config"
	"net/http"
	"time"
)

func NewServer(
	cfg *config.Server,
	router http.Handler,
) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
}
