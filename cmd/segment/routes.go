package main

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/api/middleware"
	"avito-internship-2023/internal/handler/ping"
	"avito-internship-2023/internal/handler/segment"
	"avito-internship-2023/internal/handler/user"
	"avito-internship-2023/internal/repository"
	segmentRepo "avito-internship-2023/internal/repository/segment"
	segmentSvc "avito-internship-2023/internal/service/segment"
	"context"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

func useMiddlewares(r *mux.Router) {
	r.Use(middleware.Logger)
}

func runBackgroundJobs(db *sqlx.DB, period time.Duration) chan any {
	sr := segmentRepo.NewRepository(repository.NewTransactor(db))

	ticker := time.NewTicker(period)
	done := make(chan any)
	go func() {
		defer ticker.Stop()
		defer zap.L().Info("background jobs stopped")

		zap.L().Info("background jobs started")
		for {
			select {
			case <-ticker.C:
				rowsAffected, err := sr.LeftExpiredSegments(context.Background())
				if err != nil {
					zap.L().Error("failed to left expired segments", zap.Error(err))
				}

				if rowsAffected > 0 {
					zap.L().Info("expired segments left", zap.Int("rowsAffected", rowsAffected))
				}
			case <-done:
				return
			}
		}
	}()

	return done
}

func initRoutes(r *mux.Router, db *sqlx.DB) error {
	renderer := api.NewRenderer()
	transactor := repository.NewTransactor(db)

	pingHandler := ping.NewHandler(renderer)
	pingHandler.Mount(r)

	segmentService := segmentSvc.NewService(
		segmentRepo.NewRepository(transactor),
	)
	segmentHandler := segment.NewHandler(segmentService, renderer)
	segmentHandler.Mount(r)

	userHandler := user.NewHandler(segmentService, renderer)
	userHandler.Mount(r)

	useMiddlewares(r)

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		zap.L().Info("route registered", zap.String("path", tpl), zap.Strings("methods", methods))
		return nil
	})

	return nil
}
