package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
)

type newSegmentTestCase struct {
	name                    string
	slug                    string
	autoDistributionPercent int
	pre                     func(s *SegmentRepoSuite, tc *newSegmentTestCase)
	expErr                  repository.Error
}

func (s *SegmentRepoSuite) TestRepo_NewSegment() {
	testCases := []newSegmentTestCase{
		{
			name:                    "success-with-default-percent",
			slug:                    "aboba",
			autoDistributionPercent: 0,
		},
		{
			name:                    "success-with-custom-percent",
			slug:                    "aboba",
			autoDistributionPercent: 50,
		},
		{
			name: "conflict",
			pre: func(s *SegmentRepoSuite, tc *newSegmentTestCase) {
				seg := entity.NewSegment(tc.name, tc.autoDistributionPercent)
				_, err := s.r.NewSegment(context.Background(), seg)
				s.Require().NoError(err)
			},
			expErr: repository.NewAlreadyExistsError("segment already exists"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			defer s.clean()
			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			seg := entity.NewSegment(tc.name, tc.autoDistributionPercent)
			_, err := s.r.NewSegment(context.Background(), seg)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}
		})
	}
}
