package entity

import (
	"avito-internship-2023/internal/dto"
	"github.com/google/uuid"
)

type Segment struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
}

func NewSegment(slug string) *Segment {
	return &Segment{
		ID:   uuid.New(),
		Slug: slug,
	}
}

func (s *Segment) ToSegmentCreatedDTO() *dto.SegmentCreated {
	return &dto.SegmentCreated{
		ID:   s.ID.String(),
		Slug: s.Slug,
	}
}

func (s *Segment) ToUserSegmentDTO() *dto.UserSegment {
	return &dto.UserSegment{
		Slug: s.Slug,
	}
}
