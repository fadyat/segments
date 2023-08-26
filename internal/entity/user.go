package entity

import (
	"avito-internship-2023/internal/dto"
	"github.com/google/uuid"
	"time"
)

type UserSegment struct {
	UserID   uint64     `db:"user_id"`
	SlugID   uuid.UUID  `db:"segment_id"`
	Slug     string     `db:"segment_slug"`
	JoinedAt *time.Time `db:"joined_at"`
	LeftAt   *time.Time `db:"left_at"`
	DueAt    *time.Time `db:"due_at"`
}

func (s *UserSegment) ToUserSegmentDTO() *dto.UserSegment {
	return &dto.UserSegment{
		Slug: s.Slug,
	}
}

func (s *UserSegment) ToHistoryReportLeftDTO() *dto.HistoryReport {
	return &dto.HistoryReport{
		UserID:        s.UserID,
		Slug:          s.Slug,
		HappenedAt:    s.LeftAt.Format(time.RFC3339),
		OperationType: "left",
	}
}

func (s *UserSegment) ToHistoryReportJoinedDTO() *dto.HistoryReport {
	return &dto.HistoryReport{
		UserID:        s.UserID,
		Slug:          s.Slug,
		HappenedAt:    s.JoinedAt.Format(time.RFC3339),
		OperationType: "joined",
	}
}
