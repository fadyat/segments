package segment

import (
	"avito-internship-2023/internal/repository"
	"context"
	"database/sql"
)

func (r *repo) LeftExpiredSegments(ctx context.Context) (int, error) {
	txOpts := &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}

	var left int
	e := r.RunTx(ctx, txOpts, func(ctx context.Context) (err error) {
		r.UseTx(ctx, func(e repository.Executor) {
			left, err = r.leftExpiredSegments(ctx, e)
		})
		return
	})

	return left, toSegmentError(e)
}

func (r *repo) leftExpiredSegments(
	ctx context.Context, e repository.Executor,
) (int, error) {
	query := `
		update user_segment
		set left_at = now()
		where due_at is not null and
		      due_at < now() and
		      left_at is null
	`

	res, err := e.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), nil
}
