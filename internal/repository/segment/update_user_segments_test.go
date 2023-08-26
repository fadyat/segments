package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
)

type updateUserSegmentsTestCase struct {
	name   string
	pre    func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase)
	userID uint64
	slugs  []string
	expErr error
}

func (s *SegmentRepoSuite) TestRepo_JoinUserToSegments() {
	validate := func(s *SegmentRepoSuite, userID uint64, slugs []string) {
		segments, err := s.r.GetActiveUserSegments(context.Background(), userID)
		s.Require().NoError(err)
		s.Require().Equal(len(slugs), len(segments))

		for _, segment := range segments {
			s.Require().Contains(slugs, segment.Slug)
		}
	}

	testCases := []updateUserSegmentsTestCase{
		{
			name:   "segments not found",
			userID: 1,
			slugs:  []string{"aboba", "bebra"},
			expErr: repository.NewNotFoundMultiError("invalid join segments", "aboba", "bebra"),
		},
		{
			name:   "success",
			userID: 1,
			slugs:  []string{"aboba", "bebra"},
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, slug := range tc.slugs {
					_, err := s.r.NewSegment(context.Background(), entity.NewSegment(slug))
					s.Require().NoError(err)
				}
			},
		},
		{
			name: "success with already joined",
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				_, err := s.r.NewSegment(context.Background(), entity.NewSegment("aboba"))
				s.Require().NoError(err)

				err = s.r.JoinUserToSegments(context.Background(), 1, []string{"aboba"})
				s.Require().NoError(err)
			},
			userID: 1,
			slugs:  []string{"aboba"},
		},
	}

	for _, tc := range testCases {
		if tc.pre != nil {
			tc.pre(s, &tc)
		}

		err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.slugs)
		if tc.expErr != nil {
			s.Require().Equal(tc.expErr, err)
			return
		}

		validate(s, tc.userID, tc.slugs)
		s.clean()
	}
}

func (s *SegmentRepoSuite) TestRepo_LeaveUserFromSegments() {
	validate := func(s *SegmentRepoSuite, userID uint64, slugs []string) {
		segments, err := s.r.GetActiveUserSegments(context.Background(), userID)
		s.Require().NoError(err)
		s.Require().Equal(len(slugs), len(segments))

		for _, segment := range segments {
			s.Require().Contains(slugs, segment.Slug)
		}
	}

	testCases := []updateUserSegmentsTestCase{
		{
			name:   "segments not found",
			userID: 1,
			slugs:  []string{"aboba", "bebra"},
			expErr: repository.NewNotFoundMultiError("invalid leave segments", "aboba", "bebra"),
		},
		{
			name:   "success",
			userID: 1,
			slugs:  []string{"aboba", "bebra"},
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, slug := range tc.slugs {
					_, err := s.r.NewSegment(context.Background(), entity.NewSegment(slug))
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), 1, tc.slugs)
				s.Require().NoError(err)
			},
		},
		{
			name: "success with already left",
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				_, err := s.r.NewSegment(context.Background(), entity.NewSegment("aboba"))
				s.Require().NoError(err)

				err = s.r.JoinUserToSegments(context.Background(), 1, []string{"aboba"})
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), 1, []string{"aboba"})
				s.Require().NoError(err)
			},
			userID: 1,
			slugs:  []string{"aboba"},
		},
	}

	for _, tc := range testCases {
		if tc.pre != nil {
			tc.pre(s, &tc)
		}

		err := s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.slugs)
		if tc.expErr != nil {
			s.Require().Equal(tc.expErr, err)
			return
		}

		validate(s, tc.userID, tc.slugs)
		s.clean()
	}
}
