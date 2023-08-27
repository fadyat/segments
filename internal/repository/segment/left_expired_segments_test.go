package segment

import (
	"avito-internship-2023/internal/entity"
	"context"
	"time"
)

type leftExpiredSegmentsTestcase struct {
	name         string
	pre          func(s *SegmentRepoSuite, tc *leftExpiredSegmentsTestcase)
	preSegments  []*entity.UserSegment
	rowsAffected int
	userID       uint64
}

func (s *SegmentRepoSuite) TestRepo_LeftExpiredSegmentsTest() {
	testCases := []leftExpiredSegmentsTestcase{
		{
			name: "success",
			pre: func(s *SegmentRepoSuite, tc *leftExpiredSegmentsTestcase) {
				for _, segment := range tc.preSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preSegments)
				s.Require().NoError(err)
			},
			preSegments: []*entity.UserSegment{
				{Slug: "aboba", DueAt: toTimePtr(time.Now().Add(-time.Hour).UTC())},
				{Slug: "bebra", DueAt: toTimePtr(time.Now().Add(-time.Hour).UTC())},
				{Slug: "kekis", DueAt: toTimePtr(time.Now().Add(-time.Hour).UTC())},
			},
			rowsAffected: 3,
			userID:       1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			left, err := s.r.LeftExpiredSegments(context.Background())
			s.Require().NoError(err)
			s.Require().Equal(tc.rowsAffected, left)

			segments, err := s.r.GetSegmentsByLeftAt(context.Background(), currentYearAndMonth())
			s.Require().NoError(err)
			s.Require().Equal(len(tc.preSegments), len(segments))

			for _, segment := range segments {
				now := time.Now().UTC()
				s.Require().True(equalTimes(&now, segment.LeftAt))
			}
		})
	}
}
