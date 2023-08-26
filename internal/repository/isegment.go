package repository

import (
	"avito-internship-2023/internal/entity"
	"context"
	"github.com/google/uuid"
	"time"
)

type ISegment interface {
	Transactor

	NewSegment(ctx context.Context, segment *entity.Segment) (createdSegment *entity.Segment, err error)
	DeleteSegment(ctx context.Context, id uuid.UUID) (err error)

	JoinUserToSegments(ctx context.Context, userID uint64, slugs []string) error
	LeaveUserFromSegments(ctx context.Context, userID uint64, slugs []string) error

	GetActiveUserSegments(ctx context.Context, userID uint64) (segments []*entity.UserSegment, err error)

	GetSegmentsByJoinAt(ctx context.Context, period time.Time) (segments []*entity.UserSegment, err error)
	GetSegmentsByLeftAt(ctx context.Context, period time.Time) (segments []*entity.UserSegment, err error)
}
