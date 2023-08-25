package segment

import (
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"go.uber.org/zap"
)

func (r *repo) GetUserSegments(
	ctx context.Context, status dto.UserSegmentStatus, userID uint64,
) (segments []*entity.Segment, err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		segments, err = r.getUserSegments(ctx, e, status, userID)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) getUserSegments(
	ctx context.Context, e repository.Executor, _ dto.UserSegmentStatus, userID uint64,
) ([]*entity.Segment, error) {
	selectQuery := `
		select id, slug from segment
		join user_segment us on segment.id = us.segment_id
		where us.user_id = $1
	`

	segments := make([]*entity.Segment, 0)
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
		var segment entity.Segment
		if er := row.StructScan(&segment); er != nil {
			return nil, er
		}

		segments = append(segments, &segment)
	}

	return segments, nil
}
