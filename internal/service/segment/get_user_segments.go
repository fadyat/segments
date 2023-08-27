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

func (s *svc) GetActiveUserSegments(
	ctx context.Context, userID string,
) ([]*dto.UserSegment, error) {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, api.NewBadRequestError(fmt.Sprintf("invalid id format: %s", userID))
	}

	if id <= 0 {
		return nil, api.NewBadRequestError("id must be positive")
	}

	var segments []*entity.UserSegment
	txOpts := &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}
	e := s.segmentRepository.RunTx(ctx, txOpts, func(ctx context.Context) error {
		segments, err = s.segmentRepository.GetActiveUserSegments(ctx, id)
		return err
	})

	var known repository.Error
	switch {
	case errors.As(e, &known):
		return nil, known.ToApiError()
	case e != nil:
		return nil, e
	}

	var dtoSegments = make([]*dto.UserSegment, 0, len(segments))
	for _, segment := range segments {
		dtoSegments = append(dtoSegments, segment.ToUserSegmentDTO())
	}

	return dtoSegments, nil
}
