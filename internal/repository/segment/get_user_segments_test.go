package segment

import (
	"avito-internship-2023/internal/dto"
	"avito-internship-2023/internal/entity"
	"context"
)

type getUserSegmentsTestCase struct {
	name     string
	preSlugs []string
	pre      func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase)
	status   dto.UserSegmentStatus
	userID   uint64
	expErr   error
}

func (s *SegmentRepoSuite) TestRepo_GetUserSegments() {
	testCases := []getUserSegmentsTestCase{
		{
			name:     "success",
			userID:   1,
			preSlugs: []string{"aboba", "bebra"},
			pre: func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase) {
				for _, slug := range tc.preSlugs {
					_, err := s.r.NewSegment(context.Background(), entity.NewSegment(slug))
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preSlugs)
				s.Require().NoError(err)
			},
			status: dto.Active,
		},
		{
			name:     "success zero segments",
			status:   dto.Active,
			preSlugs: make([]string, 0),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			segments, err := s.r.GetUserSegments(context.Background(), tc.status, tc.userID)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			s.Require().NoError(err)
			s.Require().Equal(len(tc.preSlugs), len(segments))
			for _, segment := range segments {
				s.Require().Contains(tc.preSlugs, segment.Slug)
			}
		})
	}
}
