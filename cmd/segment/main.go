package main

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/config"
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (

	// Version is a v*.*.* version of the application.
	// Passed during build time, e.g.:
	// go build -ldflags "-X main.Version=1.0.0"
	Version = "dev"
)

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("failed to init logger", zap.Error(err))
	}

	zap.ReplaceGlobals(logger)
}

func initDatabase() *sqlx.DB {
	databaseConfig, err := config.NewDatabase()
	if err != nil {
		zap.L().Fatal("failed to init database config", zap.Error(err))
	}

	db, err := sqlx.ConnectContext(
		context.Background(),
		databaseConfig.DriverName,
		databaseConfig.DataSourceName,
	)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}

	if e := applyMigrations(databaseConfig, db.DB); e != nil {
		zap.L().Fatal("failed to apply migrations", zap.Error(e))
	}

	zap.L().Info("database connection established, migrations applied")
	return db
}

func applyMigrations(cfg *config.Database, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, cfg.DriverName, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func initServer(router *mux.Router) *http.Server {
	serverConfig, err := config.NewServer()
	if err != nil {
		zap.L().Fatal("failed to init server config", zap.Error(err))
	}

	zap.L().Info(
		"starting server",
		zap.String("version", Version),
		zap.String("addr", serverConfig.Addr),
	)

	return api.NewServer(serverConfig, router)
}

func main() {
	initLogger()
	db := initDatabase()
	backgroundJobsDoneCh := runBackgroundJobs(db)

	v1 := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	initRoutes(v1, db)

	server := initServer(v1)
	go func() {
		if e := server.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			zap.L().Fatal("failed to start server", zap.Error(e))
		}
	}()

	stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	<-stopContext.Done()
	close(backgroundJobsDoneCh)

	zap.L().Info("stopping server")
	if e := server.Shutdown(context.Background()); e != nil {
		zap.L().Fatal("failed to stop server", zap.Error(e))
	}
}
