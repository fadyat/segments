package segment

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"errors"
)

func (s *svc) NewSegment(ctx context.Context, createSegment *dto.CreateSegment) (*dto.SegmentCreated, error) {
	err := createSegment.Validate()
	if err != nil {
		return nil, api.NewUnprocessableEntityError("validation error", err.Error())
	}

	var createdSegment *entity.Segment
	e := s.segmentRepository.RunTx(ctx, nil, func(ctx context.Context) error {
		createdSegment, err = s.segmentRepository.NewSegment(ctx, entity.NewSegment(createSegment.Slug))
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
		return nil, err
	}

	return createdSegment.ToSegmentCreatedDTO(), nil
}
