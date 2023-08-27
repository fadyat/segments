package segment

import (
	"avito-internship-2023/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

func (r *repo) JoinUsersToSegmentAuto(ctx context.Context) (int, error) {
	var joined int
	txOpts := &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}
	e := r.RunTx(ctx, txOpts, func(ctx context.Context) (err error) {
		r.UseTx(ctx, func(e repository.Executor) {
			joined, err = r.joinUsersToSegmentAuto(ctx, e)
		})
		return
	})
	return joined, toSegmentError(e)
}

type segmentWithNeededUsers struct {
	ID          string `db:"id"`
	Slug        string `db:"slug"`
	NeededUsers int    `db:"needed_users"`
}

func (r *repo) joinUsersToSegmentAuto(
	ctx context.Context, e repository.Executor,
) (int, error) {
	totalUsers, err := r.getTotalUsersCount(ctx, e)
	if err != nil {
		zap.L().Error("failed to get total users count", zap.Error(err))
		return 0, err
	}

	segments, err := r.getSegmentsWithNotEnoughUsers(ctx, e, totalUsers)
	if err != nil {
		zap.L().Error("failed to get segments with not enough users", zap.Error(err))
		return 0, err
	}

	var joined int
	for _, segment := range segments {
		rowsAffected, er := r.joinAvailableUsersToSegment(ctx, e, segment.ID, segment.NeededUsers)
		if er != nil {
			zap.L().Error("failed to join available users to segment", zap.Error(er))
			return 0, er
		}

		zap.L().Info("auto joined users to segment", zap.String("segment", segment.Slug), zap.Int("joined", rowsAffected))
		joined += rowsAffected
	}

	return joined, nil
}

func (r *repo) getTotalUsersCount(ctx context.Context, e repository.Executor) (int, error) {
	query := `select count(*) from "user"`
	var count int
	err := e.GetContext(ctx, &count, query)
	return count, err
}

func (r *repo) getSegmentsWithNotEnoughUsers(
	ctx context.Context, e repository.Executor, totalUsers int,
) ([]*segmentWithNeededUsers, error) {
	query := `
		with non_left_user_segment as (
			select *
			from user_segment us
			where us.left_at is null
		)
		select s.id,
		       s.slug,
		       ($1 * s.auto_distribution_percent) / 100 - count(us.user_id) as needed_users
		from segment s
		left join non_left_user_segment us on us.segment_id = s.id
		where s.auto_distribution_percent > 0 
		group by s.id
		having ($1 * s.auto_distribution_percent) / 100 - count(us.user_id) > 0
	`

	var segments = make([]*segmentWithNeededUsers, 0)
	err := e.SelectContext(ctx, &segments, query, totalUsers)
	return segments, err
}

func (r *repo) joinAvailableUsersToSegment(
	ctx context.Context, e repository.Executor, segmentID string, need int,
) (int, error) {
	availableUsers, err := r.getAvailableUsersForSegment(ctx, e, segmentID, need)
	if err != nil {
		return 0, err
	}

	return r.joinUsersToSegment(ctx, e, segmentID, availableUsers)
}

func (r *repo) joinUsersToSegment(
	ctx context.Context, e repository.Executor, segmentID string, users []uint64,
) (int, error) {
	var (
		queryBuilder strings.Builder
		args         = make([]any, 0, 2*len(users))
	)

	queryBuilder.WriteString("insert into user_segment (user_id, segment_id) values ")
	for i, userID := range users {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
		args = append(args, userID, segmentID)
	}

	queryBuilder.WriteString(" on conflict (user_id, segment_id) where left_at is null do nothing")
	res, err := e.ExecContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rows), nil
}

func (r *repo) getAvailableUsersForSegment(
	ctx context.Context, e repository.Executor, segmentID string, need int,
) ([]uint64, error) {
	query := `
		with used_users as (
				select segment.id, user_segment.user_id
                from user_segment
                	join segment on user_segment.segment_id = segment.id
                where segment.id = $1
		)
		select "user".id
		from "user"
		where "user".id not in (select user_id from used_users)
		order by random()
		limit $2
	`

	var users = make([]uint64, 0)
	err := e.SelectContext(ctx, &users, query, segmentID, need)
	return users, err
}
