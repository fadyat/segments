package service

import (
	"avito-internship-2023/internal/dto"
	"context"
)

type ISegment interface {
	NewSegment(context.Context, *dto.CreateSegment) (*dto.SegmentCreated, error)
	DeleteSegment(context.Context, string) error

	// UpdateUserSegments is a method that updates user segments,
	// can add user to segments and remove user from segments.
	//
	// If there's no such segment, and validation error will be returned.
	UpdateUserSegments(context.Context, string, *dto.UpdateUserSegments) error

	// GetActiveUserSegments is a method that returns all active user segments.
	GetActiveUserSegments(context.Context, string) ([]*dto.UserSegment, error)

	// GetHistoryReport is a method that returns history report for
	// year-month period.
	GetHistoryReport(ctx context.Context, period string) (*dto.HistoryReportList, error)

	// JoinSegmentsWithTTL is a method that joins user to segments with ttl.
	//
	// Example: user want to join in ABOBA segment for 10 days.
	JoinSegmentsWithTTL(context.Context, string, []*dto.SegmentWithTTL) error
}
