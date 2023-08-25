package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
)

func (s *SegmentRepoSuite) TestRepo_NewSegment() {
	testCases := []struct {
		name    string
		pre     func(s *SegmentRepoSuite)
		segment *entity.Segment
		expErr  repository.Error
	}{
		{
			name:    "success",
			segment: entity.NewSegment("aboba"),
		},
		{
			name:    "conflict",
			segment: entity.NewSegment("aboba"),
			pre: func(s *SegmentRepoSuite) {
				_, err := s.r.NewSegment(context.Background(), entity.NewSegment("aboba"))
				s.Require().NoError(err)
			},
			expErr: repository.NewAlreadyExistsError("segment already exists"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.pre != nil {
				tc.pre(s)
			}

			_, err := s.r.NewSegment(context.Background(), tc.segment)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			s.clean()
		})
	}
}
