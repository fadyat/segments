package dto

import (
	"errors"
	"strings"
)

type CreateSegment struct {
	Slug string `json:"slug"`
}

func (c *CreateSegment) Validate() error {
	var errorBuilder strings.Builder
	if len(c.Slug) < 3 {
		errorBuilder.WriteString("slug is too short, min length is 3\n")
	}

	if len(c.Slug) > 300 {
		errorBuilder.WriteString("slug is too long, max length is 300\n")
	}

	if errorBuilder.Len() > 0 {
		return errors.New(errorBuilder.String())
	}

	return nil
}

type SegmentCreated struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}
