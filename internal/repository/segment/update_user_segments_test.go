package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"time"
)

type updateUserSegmentsTestCase struct {
	name         string
	pre          func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase)
	userID       uint64
	userSegments []*entity.UserSegment
	expErr       error
}

// containsUserSegments may work incorrectly when not all dates are provided.
//
// in context of current tests everything is ok.
func containsUserSegments(segments []*entity.UserSegment, segment *entity.UserSegment) bool {
	for _, s := range segments {
		equalSlugs := s.Slug == segment.Slug
		equalDueAt := equalTimes(s.DueAt, segment.DueAt)
		equalLeftAt := equalTimes(s.LeftAt, segment.LeftAt)

		if equalSlugs && equalDueAt && equalLeftAt {
			return true
		}
	}

	return false
}

func validateUserSegments(
	s *SegmentRepoSuite, userID uint64, userSegments []*entity.UserSegment,
) {
	s.T().Helper()

	segments, err := s.r.GetActiveUserSegments(context.Background(), userID)
	s.Require().NoError(err)
	s.Require().Equal(len(userSegments), len(segments))

	for _, segment := range segments {
		s.Require().True(containsUserSegments(userSegments, segment), "segment not found")
	}
}

func (s *SegmentRepoSuite) TestRepo_JoinUserToSegments() {
	testCases := []updateUserSegmentsTestCase{
		{
			name:   "segments not found",
			userID: 1,
			userSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
			},
			expErr: repository.NewNotFoundMultiError("invalid join segments", "aboba", "bebra"),
		},
		{
			name:   "success",
			userID: 1,
			userSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
			},
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, us := range tc.userSegments {
					sg := entity.NewSegment(us.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}
			},
		},
		{
			name: "success with already joined",
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, us := range tc.userSegments {
					sg := entity.NewSegment(us.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.userSegments)
				s.Require().NoError(err)
			},
			userID: 1,
			userSegments: []*entity.UserSegment{
				{Slug: "aboba"},
			},
		},
		{
			name: "success with due at",
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, us := range tc.userSegments {
					sg := entity.NewSegment(us.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}
			},
			userID: 1,
			userSegments: []*entity.UserSegment{
				{Slug: "aboba", DueAt: toTimePtr(time.Now().Add(time.Hour))},
				{Slug: "bebra", DueAt: toTimePtr(time.Now().Add(time.Minute))},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			defer s.clean()

			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.userSegments)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			validateUserSegments(s, tc.userID, tc.userSegments)
		})
	}
}

func (s *SegmentRepoSuite) TestRepo_LeaveUserFromSegments() {
	testCases := []updateUserSegmentsTestCase{
		{
			name:   "segments not found",
			userID: 1,
			userSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
			},
			expErr: repository.NewNotFoundMultiError("invalid leave segments", "aboba", "bebra"),
		},
		{
			name:   "success",
			userID: 1,
			userSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
			},
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, segment := range tc.userSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.userSegments)
				s.Require().NoError(err)
			},
		},
		{
			name: "success with already left",
			userSegments: []*entity.UserSegment{
				{Slug: "aboba"},
			},
			pre: func(s *SegmentRepoSuite, tc *updateUserSegmentsTestCase) {
				for _, segment := range tc.userSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.userSegments)
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.userSegments)
				s.Require().NoError(err)
			},
			userID: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			defer s.clean()

			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			err := s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.userSegments)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			validateUserSegments(s, tc.userID, nil)
		})
	}
}
