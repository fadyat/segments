package service

import (
	"avito-internship-2023/internal/dto"
	"context"
)

type ISegment interface {
	NewSegment(context.Context, *dto.CreateSegment) (*dto.SegmentCreated, error)
	DeleteSegment(context.Context, string) error
}
