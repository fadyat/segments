package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func (s *svc) DeleteSegment(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return api.NewBadRequestError(fmt.Sprintf("invalid id format: %s", id))
	}

	e := s.segmentRepository.RunTx(ctx, nil, func(ctx context.Context) error {
		return s.segmentRepository.DeleteSegment(ctx, uid)
	})

	var known repository.Error
	switch {
	case errors.As(e, &known):
		return known.ToApiError()
	case e != nil:
		return e
	}

	return nil
}
