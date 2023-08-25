package repository

import (
	"avito-internship-2023/internal/entity"
	"context"
)

type ISegment interface {
	Transactor

	NewSegment(ctx context.Context, segment *entity.Segment) (createdSegment *entity.Segment, err error)
}
