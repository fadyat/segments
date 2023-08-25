package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
)

func (r *repo) NewSegment(
	ctx context.Context, segment *entity.Segment,
) (createdSegment *entity.Segment, err error) {
	r.UseTx(ctx, func(executor repository.Executor) {
		createdSegment, err = r.newSegment(ctx, executor, segment)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) newSegment(
	ctx context.Context, executor repository.Executor, segment *entity.Segment,
) (*entity.Segment, error) {
	insertQuery := `
		insert into segment (id, slug)
		values ($1, $2)`

	_, err := executor.ExecContext(ctx, insertQuery, segment.ID, segment.Slug)
	if err != nil {
		return nil, err
	}

	return segment, nil
}
