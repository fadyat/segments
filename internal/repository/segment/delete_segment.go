package segment

import (
	"avito-internship-2023/internal/repository"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

func (r *repo) DeleteSegment(ctx context.Context, id uuid.UUID) (err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		err = r.deleteSegment(ctx, e, id)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) deleteSegment(
	ctx context.Context, executor repository.Executor, id uuid.UUID,
) error {
	deleteUserSegmentQuery := `
		delete from user_segment
		where segment_id = $1
	`

	_, err := executor.ExecContext(ctx, deleteUserSegmentQuery, id)
	if err != nil {
		return err
	}

	deleteQuery := `
		delete from segment
		where id = $1	
	`

	result, err := executor.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
