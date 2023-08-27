package segment

import (
	"avito-internship-2023/internal/entity"
	"context"
	"fmt"
	"time"
)

type getHistoryReportTestCase struct {
	name             string
	preUserSegments  []*entity.UserSegment
	expectedSegments []*entity.UserSegment
	pre              func(s *SegmentRepoSuite, tc *getHistoryReportTestCase)
	userID           uint64
	period           time.Time
}

func currentYearAndMonth() time.Time {
	t, _ := time.Parse("2006-01", time.Now().Format("2006-01"))
	return t
}

func (s *SegmentRepoSuite) TestRepo_GetSegmentsByJoinAt() {
	testCases := []getHistoryReportTestCase{
		{
			name: "success",
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis"},
			},
			expectedSegments: []*entity.UserSegment{
				{Slug: "aboba", JoinedAt: toTimePtr(time.Now())},
				{Slug: "bebra", JoinedAt: toTimePtr(time.Now())},
				{Slug: "kekis", JoinedAt: toTimePtr(time.Now())},
			},
			pre: func(s *SegmentRepoSuite, tc *getHistoryReportTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)
			},
			userID: 1,
			period: currentYearAndMonth(),
		},
		{
			name: "not in period",
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis"},
			},
			pre: func(s *SegmentRepoSuite, tc *getHistoryReportTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)
			},
			userID:           1,
			period:           currentYearAndMonth().AddDate(0, -1, 0),
			expectedSegments: []*entity.UserSegment{},
		},
		{
			name: "one in period",
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis", JoinedAt: toTimePtr(time.Now().AddDate(0, -1, 0).UTC())},
			},
			pre: func(s *SegmentRepoSuite, tc *getHistoryReportTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)
			},
			userID: 1,
			period: currentYearAndMonth().AddDate(0, -1, 0),
			expectedSegments: []*entity.UserSegment{
				{Slug: "kekis", JoinedAt: toTimePtr(time.Now().AddDate(0, -1, 0))},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			defer s.clean()

			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			segments, err := s.r.GetSegmentsByJoinAt(context.Background(), tc.period)
			s.Require().NoError(err)
			s.Require().Equal(len(tc.expectedSegments), len(segments))

			for i := range segments {
				s.Require().Equal(tc.expectedSegments[i].Slug, segments[i].Slug)

				expUTC := tc.expectedSegments[i].JoinedAt.UTC()
				actUTC := segments[i].JoinedAt.UTC()
				msg := fmt.Sprintf("expected: %s, actual: %s", expUTC, actUTC)
				s.Require().True(equalTimes(&expUTC, &actUTC), msg)
			}
		})
	}
}

func (s *SegmentRepoSuite) TestRepo_GetSegmentsByLeftAt() {
	testCases := []getHistoryReportTestCase{
		{
			name: "success",
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis"},
			},
			expectedSegments: []*entity.UserSegment{
				{Slug: "aboba", LeftAt: toTimePtr(time.Now())},
				{Slug: "bebra", LeftAt: toTimePtr(time.Now())},
				{Slug: "kekis", LeftAt: toTimePtr(time.Now())},
			},
			pre: func(s *SegmentRepoSuite, tc *getHistoryReportTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)
			},
			userID: 1,
			period: currentYearAndMonth(),
		},
		{
			name: "not in period",
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis"},
			},
			pre: func(s *SegmentRepoSuite, tc *getHistoryReportTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)
			},
			userID:           1,
			period:           currentYearAndMonth().AddDate(0, -1, 0),
			expectedSegments: []*entity.UserSegment{},
		},
		{
			name: "one in period",
			preUserSegments: []*entity.UserSegment{
				{Slug: "aboba"},
				{Slug: "bebra"},
				{Slug: "kekis", LeftAt: toTimePtr(time.Now().AddDate(0, -1, 0).UTC())},
			},
			pre: func(s *SegmentRepoSuite, tc *getHistoryReportTestCase) {
				for _, segment := range tc.preUserSegments {
					sg := entity.NewSegment(segment.Slug, 0)
					_, err := s.r.NewSegment(context.Background(), sg)
					s.Require().NoError(err)
				}

				err := s.r.JoinUserToSegments(context.Background(), tc.userID, tc.preUserSegments)
				s.Require().NoError(err)

				err = s.r.LeaveUserFromSegments(context.Background(), tc.userID, tc.preUserSegments)
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

			segments, err := s.r.GetSegmentsByLeftAt(context.Background(), tc.period)
			s.Require().NoError(err)
			s.Require().Equal(len(tc.expectedSegments), len(segments))

			for i := range segments {
				s.Require().Equal(tc.expectedSegments[i].Slug, segments[i].Slug)

				expUTC := tc.expectedSegments[i].LeftAt.UTC()
				actUTC := segments[i].LeftAt.UTC()
				msg := fmt.Sprintf("expected: %s, actual: %s", expUTC, actUTC)
				s.Require().True(equalTimes(&expUTC, &actUTC), msg)
			}
		})
	}
}
