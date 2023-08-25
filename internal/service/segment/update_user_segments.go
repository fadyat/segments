package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
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
	e := s.segmentRepository.RunTransaction(ctx, txOpts, func(ctx context.Context) error {
		if updateUserSegmentsDTO.CanJoin() {
			err = s.segmentRepository.JoinUserToSegments(ctx, id, updateUserSegmentsDTO.JoinSegments)
			if err != nil {
				return err
			}
		}

		if updateUserSegmentsDTO.CanLeave() {
			err = s.segmentRepository.LeaveUserFromSegments(ctx, id, updateUserSegmentsDTO.LeaveSegments)
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
