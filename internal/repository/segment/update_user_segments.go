package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"fmt"
	"strings"
)

func (r *repo) JoinUserToSegments(
	ctx context.Context, userID uint64, segments []*entity.UserSegment,
) (err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		err = r.joinUserToSegments(ctx, e, userID, segments)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) joinUserToSegments(
	ctx context.Context, e repository.Executor, userID uint64, segments []*entity.UserSegment,
) (err error) {
	knownSegments, err := r.getKnownSegmentsBySlugs(ctx, e, segments)
	if err != nil {
		return err
	}

	unknownSlugs := r.getUnknownSlugs(knownSegments, segments)
	if len(unknownSlugs) > 0 {
		return repository.NewNotFoundMultiError("invalid join segments", unknownSlugs...)
	}

	var (
		queryBuilder strings.Builder
		args         = make([]any, 0, 4*len(knownSegments))
	)

	queryBuilder.WriteString("insert into user_segment (user_id, segment_id, due_at, joined_at) values ")
	for i, segment := range knownSegments {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		args = append(args, userID, segment.ID, segments[i].DueAt, segments[i].JoinedAt)
	}

	queryBuilder.WriteString(" on conflict (user_id, segment_id) where left_at is null do nothing")
	_, err = e.ExecContext(ctx, queryBuilder.String(), args...)
	return err
}

func (r *repo) LeaveUserFromSegments(
	ctx context.Context, userID uint64, segments []*entity.UserSegment,
) (err error) {
	r.UseTx(ctx, func(e repository.Executor) {
		err = r.leaveUserFromSegments(ctx, e, userID, segments)
		err = toSegmentError(err)
	})
	return
}

func (r *repo) leaveUserFromSegments(
	ctx context.Context, e repository.Executor, userID uint64, segments []*entity.UserSegment,
) (err error) {
	knownSegments, err := r.getKnownSegmentsBySlugs(ctx, e, segments)
	if err != nil {
		return err
	}

	unknownSlugs := r.getUnknownSlugs(knownSegments, segments)
	if len(unknownSlugs) > 0 {
		return repository.NewNotFoundMultiError("invalid leave segments", unknownSlugs...)
	}

	var (
		queryBuilder strings.Builder
		args         = make([]any, 0, 2*len(knownSegments))
	)
	queryBuilder.WriteString("update user_segment set left_at = now() where (user_id, segment_id) in (")
	for i, segment := range knownSegments {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
		args = append(args, userID, segment.ID)
	}

	queryBuilder.WriteString(") and left_at is null")
	_, err = e.ExecContext(ctx, queryBuilder.String(), args...)
	return err
}
