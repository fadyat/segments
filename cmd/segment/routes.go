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
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func useMiddlewares(r *mux.Router) {
	r.Use(middleware.Logger)
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
