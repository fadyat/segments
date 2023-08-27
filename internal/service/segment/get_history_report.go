package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"time"
)

func (s *svc) GetHistoryReport(ctx context.Context, period string) (*dto.HistoryReportList, error) {
	timeRange, err := time.Parse("2006-01", period)
	if err != nil {
		zap.L().Debug("failed to parse time range", zap.Error(err))
		return nil, api.NewBadRequestError("invalid time range format, expected YYYY-MM")
	}

	var (
		joinAtSegments []*entity.UserSegment
		leftAtSegments []*entity.UserSegment
	)
	txOpts := &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true}
	e := s.segmentRepository.RunTx(ctx, txOpts, func(ctx context.Context) error {
		joinAtSegments, err = s.segmentRepository.GetSegmentsByJoinAt(ctx, timeRange)
		if err != nil {
			return err
		}

		leftAtSegments, err = s.segmentRepository.GetSegmentsByLeftAt(ctx, timeRange)
		if err != nil {
			return err
		}

		return nil
	})

	var known repository.Error
	switch {
	case errors.As(e, &known):
		return nil, known.ToApiError()
	case e != nil:
		return nil, e
	}

	dtoSegments := make([]*dto.HistoryReport, 0, len(joinAtSegments)+len(leftAtSegments))
	var joinAtIndex, leftAtIndex int
	for joinAtIndex < len(joinAtSegments) || leftAtIndex < len(leftAtSegments) {
		if joinAtIndex == len(joinAtSegments) {
			dtoSegments = append(dtoSegments, leftAtSegments[leftAtIndex].ToHistoryReportLeftDTO())
			leftAtIndex++
			continue
		}

		if leftAtIndex == len(leftAtSegments) {
			dtoSegments = append(dtoSegments, joinAtSegments[joinAtIndex].ToHistoryReportJoinedDTO())
			joinAtIndex++
			continue
		}

		joinedAt := joinAtSegments[joinAtIndex].JoinedAt
		leftAt := leftAtSegments[leftAtIndex].LeftAt
		if joinedAt.Before(*leftAt) {
			dtoSegments = append(dtoSegments, joinAtSegments[joinAtIndex].ToHistoryReportJoinedDTO())
			joinAtIndex++
		} else {
			dtoSegments = append(dtoSegments, leftAtSegments[leftAtIndex].ToHistoryReportLeftDTO())
			leftAtIndex++
		}
	}

	return &dto.HistoryReportList{Segments: dtoSegments}, nil
}
