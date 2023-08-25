package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"fmt"
	"strings"
)

func (r *repo) getKnownSegmentsBySlugs(
	ctx context.Context, e repository.Executor, slugs []string,
) ([]*entity.Segment, error) {
	var (
		queryBuilder strings.Builder
		args         = make([]any, 0, len(slugs))
	)

	queryBuilder.WriteString("select * from segment where slug in (")
	for i, slug := range slugs {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("$%d", i+1))
		args = append(args, slug)
	}

	queryBuilder.WriteString(")")

	var segments = make([]*entity.Segment, 0, len(slugs))
	err := e.SelectContext(ctx, &segments, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func (r *repo) getUnknownSlugs(knownSegments []*entity.Segment, slugs []string) []string {
	knownSlugsMap := make(map[string]bool, len(knownSegments))
	for _, segment := range knownSegments {
		knownSlugsMap[segment.Slug] = true
	}

	var unknownSlugs = make([]string, 0)
	for _, slug := range slugs {
		if _, ok := knownSlugsMap[slug]; !ok {
			unknownSlugs = append(unknownSlugs, slug)
		}
	}

	return unknownSlugs
}
