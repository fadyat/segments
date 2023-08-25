package repository

import (
	"avito-internship-2023/internal/entity"
	"context"
	"github.com/google/uuid"
)

type ISegment interface {
	Transactor

	NewSegment(ctx context.Context, segment *entity.Segment) (createdSegment *entity.Segment, err error)
	DeleteSegment(ctx context.Context, id uuid.UUID) (err error)
}
