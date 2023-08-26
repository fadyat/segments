package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"fmt"
	"strings"
)

func (r *repo) getKnownSegmentsBySlugs(
	ctx context.Context, e repository.Executor, segments []*entity.UserSegment,
) ([]*entity.Segment, error) {
	var (
		queryBuilder strings.Builder
		args         = make([]any, 0, len(segments))
	)

	queryBuilder.WriteString("select id, slug from segment where slug in (")
	for i, segment := range segments {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("$%d", i+1))
		args = append(args, segment.Slug)
	}

	queryBuilder.WriteString(")")

	var knownSegments = make([]*entity.Segment, 0, len(segments))
	err := e.SelectContext(ctx, &knownSegments, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}

	return knownSegments, nil
}

func (r *repo) getUnknownSlugs(knownSegments []*entity.Segment, segments []*entity.UserSegment) []string {
	knownSlugsMap := make(map[string]bool, len(knownSegments))
	for _, segment := range knownSegments {
		knownSlugsMap[segment.Slug] = true
	}

	var unknownSlugs = make([]string, 0)
	for _, segment := range segments {
		if _, ok := knownSlugsMap[segment.Slug]; !ok {
			unknownSlugs = append(unknownSlugs, segment.Slug)
		}
	}

	return unknownSlugs
}
