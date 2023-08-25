package segment

import (
	"avito-internship-2023/internal/entity"
	"avito-internship-2023/internal/repository"
	"context"
	"github.com/google/uuid"
)

func (s *SegmentRepoSuite) TestRepo_DeleteSegment() {
	testCases := []struct {
		name   string
		pre    func(s *SegmentRepoSuite)
		id     uuid.UUID
		expErr repository.Error
	}{
		{
			name: "success",
			pre: func(s *SegmentRepoSuite) {
				segment := entity.NewSegment("aboba")
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
			if tc.pre != nil {
				tc.pre(s)
			}

			err := s.r.DeleteSegment(context.Background(), tc.id)
			if tc.expErr != nil {
				s.Require().Equal(tc.expErr, err)
				return
			}

			s.clean()
		})
	}
}
