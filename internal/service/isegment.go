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

	// GetUserSegments is a method that returns user segments
	// filtered by status and sorted by slug.
	GetUserSegments(context.Context, dto.UserSegmentStatus, string) ([]*dto.UserSegment, error)
}
