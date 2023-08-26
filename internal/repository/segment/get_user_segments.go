package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"go.uber.org/zap"
)

func (r *repo) GetActiveUserSegments(
	ctx context.Context, userID uint64,
) (segments []*entity.UserSegment, err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		segments, err = r.getUserSegments(ctx, e, userID)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) getUserSegments(
	ctx context.Context, e repository.Executor, userID uint64,
) ([]*entity.UserSegment, error) {
	selectQuery := `
		select segment_id as segment_id,
		       slug as segment_slug,
		       user_id as user_id,
		       left_at as left_at,
		       joined_at as joined_at
		from segment s
		join user_segment us on s.id = us.segment_id
		where us.user_id = $1 and us.left_at is null
	`

	segments := make([]*entity.UserSegment, 0)
	row, err := e.QueryxContext(ctx, selectQuery, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if er := row.Close(); er != nil {
			zap.L().Error("failed to close row", zap.Error(er))
		}
	}()

	for row.Next() {
		var segment entity.UserSegment
		if er := row.StructScan(&segment); er != nil {
			return nil, er
		}

		segments = append(segments, &segment)
	}

	return segments, nil
}
