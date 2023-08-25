package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

func (s *svc) GetUserSegments(
	ctx context.Context, status dto.UserSegmentStatus, userID string,
) ([]*dto.UserSegment, error) {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, api.NewBadRequestError(fmt.Sprintf("invalid id format: %s", userID))
	}

	if id <= 0 {
		return nil, api.NewBadRequestError("id must be positive")
	}

	var segments []*entity.Segment
	txOpts := &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}
	e := s.segmentRepository.RunTransaction(ctx, txOpts, func(ctx context.Context) error {
		segments, err = s.segmentRepository.GetUserSegments(ctx, status, id)
		return err
	})

	var known repository.Error
	switch {
	case errors.As(e, &known):
		return nil, known.ToApiError()
	case e != nil:
		return nil, err
	}

	var dtoSegments = make([]*dto.UserSegment, 0, len(segments))
	for _, segment := range segments {
		dtoSegments = append(dtoSegments, segment.ToUserSegmentDTO())
	}

	return dtoSegments, nil
}
