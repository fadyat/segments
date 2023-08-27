package segment

import (
	"avito-internship-2023/internal/entity"
	"context"
)

type getUserSegmentsTestCase struct {
	name            string
	preUserSegments []*entity.UserSegment
	leaveSegments   []*entity.UserSegment
	pre             func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase)
	userID          uint64
	expErr          error
}

func (s *SegmentRepoSuite) TestRepo_GetUserSegments() {
	testCases := []getUserSegmentsTestCase{
		{
			name:   "success",
			userID: 1,
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
			},
			pre: func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)
			},
		},
		{
			name:            "success zero segments",
			preUserSegments: make([]*entity.UserSegment, 0),
		},
		{
			name:   "have left segments",
			userID: 1,
			pre: func(s *SegmentRepoSuite, tc *getUserSegmentsTestCase) {
				for _, slug := range tc.preUserSegments {
					sg := entity.NewSegment(slug.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.leaveSegments)
				s.Require().NoError(err)
			},
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis"},
			},
			leaveSegments: []*entity.UserSegment{
				{Slug: "aboba"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			defer s.clean()

			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			segments, err := s.r.GetActiveUserSegments(context.Background(), tc.userID)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			activeSegments := make([]*entity.UserSegment, 0)
			for _, segment := range tc.preUserSegments {
				if !containsUserSegments(tc.leaveSegments, segment) {
					activeSegments = append(activeSegments, segment)
				}
			}

			s.Require().NoError(err)
			s.Require().Equal(len(activeSegments), len(segments))
			for i := range activeSegments {
				s.Require().Equal(activeSegments[i].Slug, segments[i].Slug)
			}
		})
	}
}
