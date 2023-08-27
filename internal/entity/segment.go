package entity

import (
	"avito-internship-2023/internal/dto"
	"github.com/google/uuid"
)

type Segment struct {
	ID                      uuid.UUID `db:"id"`
	Slug                    string    `db:"slug"`
	AutoDistributionPercent int       `db:"auto_distribution_percent"`
}

func NewSegment(slug string, autoDistributionPercent int) *Segment {
	return &Segment{
		ID:                      uuid.New(),
		Slug:                    slug,
		AutoDistributionPercent: autoDistributionPercent,
	}
}

func (s *Segment) ToSegmentCreatedDTO() *dto.SegmentCreated {
	return &dto.SegmentCreated{
		ID:                      s.ID.String(),
		Slug:                    s.Slug,
		AutoDistributionPercent: s.AutoDistributionPercent,
	}
}
