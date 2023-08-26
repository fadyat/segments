package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (s *svc) JoinSegmentsWithTTL(
	ctx context.Context, userID string, segments []*dto.SegmentWithTTL,
) error {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return api.NewBadRequestError(fmt.Sprintf("invalid user id: %s", userID))
	}

	if id <= 0 {
		return api.NewBadRequestError("id must be positive")
	}

	e := s.segmentRepository.RunTx(ctx, nil, func(ctx context.Context) error {
		if len(segments) == 0 {
			return nil
		}

		var joinSegments = make([]*entity.UserSegment, 0, len(segments))
		for _, segment := range segments {
			dueAt := time.Now().AddDate(0, 0, segment.TTL)
			joinSegments = append(
				joinSegments,
				&entity.UserSegment{Slug: segment.Slug, DueAt: &dueAt},
			)
		}

		err = s.segmentRepository.JoinUserToSegments(ctx, id, joinSegments)
		if err != nil {
			return err
		}

		return nil
	})

	var known repository.Error
	switch {
	case errors.As(e, &known):
		return known.ToApiError()
	case e != nil:
		return err
	}

	return nil
}
