package repository

import (
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"context"
	"github.com/google/uuid"
)

type ISegment interface {
	Transactor

	NewSegment(ctx context.Context, segment *entity.Segment) (createdSegment *entity.Segment, err error)
	DeleteSegment(ctx context.Context, id uuid.UUID) (err error)
	JoinUserToSegments(ctx context.Context, userID uint64, slugs []string) error
	LeaveUserFromSegments(ctx context.Context, userID uint64, slugs []string) error
	GetUserSegments(ctx context.Context, status dto.UserSegmentStatus, userID uint64) (segments []*entity.Segment, err error)
}
