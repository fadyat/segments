package main

import (
	"avito-internship-2023/internal/config"
	"avito-internship-2023/internal/repository"
	segmentRepo "avito-internship-2023/internal/repository/segment"
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

func runBackgroundJobs(db *sqlx.DB) chan any {
	jobsConfig, err := config.NewJobs()
	if err != nil {
		zap.L().Fatal("failed to init jobs config", zap.Error(err))
	}

	sr := segmentRepo.NewRepository(repository.NewTransactor(db))

	expireTicker := time.NewTicker(jobsConfig.ExpirePeriod)
	autoJoinTicker := time.NewTicker(jobsConfig.AutoJoinPeriod)

	done := make(chan any)
	go func() {
		defer expireTicker.Stop()
		defer zap.L().Info("background jobs stopped")

		zap.L().Info("background jobs started")
		for {
			select {
			case <-expireTicker.C:
				rowsAffected, e := sr.LeftExpiredSegments(context.Background())
				if e != nil {
					zap.L().Error("failed to left expired segments", zap.Error(e))
				}

				if rowsAffected > 0 {
					zap.L().Info("expired segments left", zap.Int("rowsAffected", rowsAffected))
				}
			case <-autoJoinTicker.C:
				rowsAffected, e := sr.JoinUsersToSegmentAuto(context.Background())
				if e != nil {
					zap.L().Error("failed to join users to segment auto", zap.Error(e))
				}

				if rowsAffected > 0 {
					zap.L().Info("auto joined users to segment", zap.Int("rowsAffected", rowsAffected))
				}
			case <-done:
				return
			}
		}
	}()

	return done
}
