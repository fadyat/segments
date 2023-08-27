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

	JoinUserToSegments(ctx context.Context, userID uint64, segments []*entity.UserSegment) error
	LeaveUserFromSegments(ctx context.Context, userID uint64, segments []*entity.UserSegment) error

	// JoinUsersToSegmentAuto joins user to segment automatically via percentage of allowed users
	// to a segment.
	//
	// 0   = no one will be joined to a segment
	// 100 = all users will be joined to a segment
	//
	// Returns number of segments that were joined.
	JoinUsersToSegmentAuto(ctx context.Context) (int, error)

	// LeftExpiredSegments sets left_at to now() for all segments where due_at < now() and left_at is null.
	// Returns number of segments that were left.
	LeftExpiredSegments(context.Context) (int, error)

	GetActiveUserSegments(ctx context.Context, userID uint64) (segments []*entity.UserSegment, err error)

	GetSegmentsByJoinAt(ctx context.Context, period time.Time) (segments []*entity.UserSegment, err error)
	GetSegmentsByLeftAt(ctx context.Context, period time.Time) (segments []*entity.UserSegment, err error)
}
