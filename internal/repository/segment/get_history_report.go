package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"go.uber.org/zap"
	"time"
)

func (r *repo) GetSegmentsByJoinAt(
	ctx context.Context, period time.Time,
) (segments []*entity.UserSegment, err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		segments, err = r.getSegmentsByJoinAt(ctx, e, period)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) getSegmentsByJoinAt(
	ctx context.Context, e repository.Executor, period time.Time,
) ([]*entity.UserSegment, error) {
	selectQuery := `
		select segment_id as segment_id,
		       slug as segment_slug,
		       user_id as user_id,
		       left_at as left_at,
		       joined_at as joined_at
		from segment s	
		join user_segment us on s.id = us.segment_id
		where us.joined_at >= $1 and us.joined_at < $2
		order by us.joined_at
	`

	segments := make([]*entity.UserSegment, 0)
	rows, err := e.QueryxContext(ctx, selectQuery, period, period.AddDate(0, 1, 0))
	if err != nil {
		return nil, err
	}

	defer func() {
		if er := rows.Close(); er != nil {
			zap.L().Error("failed to close row", zap.Error(er))
		}
	}()

	for rows.Next() {
		var segment entity.UserSegment
		if er := rows.StructScan(&segment); er != nil {
			return nil, er
		}

		segments = append(segments, &segment)
	}

	return segments, nil
}

func (r *repo) GetSegmentsByLeftAt(
	ctx context.Context, period time.Time,
) (segments []*entity.UserSegment, err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		segments, err = r.getSegmentsByLeftAt(ctx, e, period)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) getSegmentsByLeftAt(
	ctx context.Context, e repository.Executor, period time.Time,
) ([]*entity.UserSegment, error) {
	selectQuery := `
		select segment_id as segment_id,
		       slug as segment_slug,
		       user_id as user_id,
		       left_at as left_at,
		       joined_at as joined_at
		from segment s
		join user_segment us on s.id = us.segment_id
		where us.left_at >= $1 and us.left_at < $2
		order by us.left_at
	`

	segments := make([]*entity.UserSegment, 0)
	rows, err := e.QueryxContext(ctx, selectQuery, period, period.AddDate(0, 1, 0))
	if err != nil {
		return nil, err
	}

	defer func() {
		if er := rows.Close(); er != nil {
			zap.L().Error("failed to close row", zap.Error(er))
		}
	}()

	for rows.Next() {
		var segment entity.UserSegment
		if er := rows.StructScan(&segment); er != nil {
			return nil, er
		}

		segments = append(segments, &segment)
	}

	return segments, nil
}
