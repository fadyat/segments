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

func (s *svc) UpdateUserSegments(
	ctx context.Context, userID string, updateUserSegmentsDTO *dto.UpdateUserSegments,
) error {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return api.NewBadRequestError(fmt.Sprintf("invalid id format: %s", userID))
	}

	if id <= 0 {
		return api.NewBadRequestError("id must be positive")
	}

	if e := updateUserSegmentsDTO.Validate(); e != nil {
		return api.NewUnprocessableEntityError("validation error", e.Error())
	}

	txOpts := &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}
	e := s.segmentRepository.RunTx(ctx, txOpts, func(ctx context.Context) error {
		if updateUserSegmentsDTO.CanJoin() {
			var joinSegments = make([]*entity.UserSegment, 0, len(updateUserSegmentsDTO.JoinSegments))
			for _, slug := range updateUserSegmentsDTO.JoinSegments {
				joinSegments = append(joinSegments, &entity.UserSegment{Slug: slug, DueAt: nil})
			}

			err = s.segmentRepository.JoinUserToSegments(ctx, id, joinSegments)
			if err != nil {
				return err
			}
		}

		if updateUserSegmentsDTO.CanLeave() {
			var leaveSegments = make([]*entity.UserSegment, 0, len(updateUserSegmentsDTO.LeaveSegments))
			for _, slug := range updateUserSegmentsDTO.LeaveSegments {
				leaveSegments = append(leaveSegments, &entity.UserSegment{Slug: slug})
			}

			err = s.segmentRepository.LeaveUserFromSegments(ctx, id, leaveSegments)
			if err != nil {
				return err
			}
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
