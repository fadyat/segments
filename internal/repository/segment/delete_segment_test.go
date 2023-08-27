package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"github.com/google/uuid"
)

type deleteSegmentTestCase struct {
	name                    string
	autoDistributionPercent int
	pre                     func(s *SegmentRepoSuite, tc *deleteSegmentTestCase)
	id                      uuid.UUID
	expErr                  repository.Error
}

func (s *SegmentRepoSuite) TestRepo_DeleteSegment() {
	testCases := []deleteSegmentTestCase{
		{
			name: "success",
			pre: func(s *SegmentRepoSuite, tc *deleteSegmentTestCase) {
				segment := entity.NewSegment(tc.name, tc.autoDistributionPercent)
				segment.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

				_, err := s.r.NewSegment(context.Background(), segment)
				s.Require().NoError(err)
			},
			id: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
		{
			name:   "not found",
			id:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			expErr: repository.NewNotFoundError("segment not found"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			defer s.clean()

			if tc.pre != nil {
				tc.pre(s, &tc)
			}

			err := s.r.DeleteSegment(context.Background(), tc.id)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}
		})
	}
}
