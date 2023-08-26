package segment

import (
	"avito-internship-2023/internal/entity"
	"context"
)

type getUserSegmentsTestCase struct {
	name       string
	preSlugs   []string
	leaveSlugs []string
	pre        func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase)
	userID     uint64
	expErr     error
}

func contains(slugs []string, slug string) bool {
	for _, s := range slugs {
		if s == slug {
			return true
		}
	}

	return false
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
		},
		{
			name:     "success zero segments",
			preSlugs: make([]string, 0),
		},
		{
			name: "have left segments",
			pre: func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase) {
				for _, slug := range tc.preSlugs {
					_, err := s.r.NewSegment(context.Background(), entity.NewSegment(slug))
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preSlugs)
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.leaveSlugs)
				s.Require().NoError(err)
			},
			preSlugs:   []string{"aboba", "bebra", "kekis"},
			leaveSlugs: []string{"aboba"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			segments, err := s.r.GetActiveUserSegments(context.Background(), tc.userID)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			activeSlugs := make([]string, 0, len(tc.preSlugs))
			for _, slug := range tc.preSlugs {
				if !contains(tc.leaveSlugs, slug) {
					activeSlugs = append(activeSlugs, slug)
				}
			}

			s.Require().NoError(err)
			s.Require().Equal(len(activeSlugs), len(segments))
			for _, segment := range segments {
				s.Require().Contains(activeSlugs, segment.Slug)
			}

			s.clean()
		})
	}
}
