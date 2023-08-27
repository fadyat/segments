package dto

import (
	"errors"
	"strings"
)

type CreateSegment struct {
	Slug                    string `json:"slug"`
	AutoDistributionPercent int    `json:"auto_distribution_percent"`
}

func (c *CreateSegment) Validate() error {
	var errorBuilder strings.Builder
	if len(c.Slug) < 3 {
		errorBuilder.WriteString("slug is too short, min length is 3\n")
	}

	if len(c.Slug) > 300 {
		errorBuilder.WriteString("slug is too long, max length is 300\n")
	}

	if c.AutoDistributionPercent < 0 {
		errorBuilder.WriteString("auto_distribution_percent must be greater than 0\n")
	}

	if c.AutoDistributionPercent > 100 {
		errorBuilder.WriteString("auto_distribution_percent must be less than 100\n")
	}

	if errorBuilder.Len() > 0 {
		return errors.New(errorBuilder.String())
	}

	return nil
}

type SegmentCreated struct {
	ID                      string `json:"id"`
	Slug                    string `json:"slug"`
	AutoDistributionPercent int    `json:"auto_distribution_percent"`
}
