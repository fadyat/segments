package dto

import "fmt"

type CreateSegment struct {
	Slug string `json:"slug"`
}

func (c *CreateSegment) Validate() error {
	if len(c.Slug) < 3 {
		return fmt.Errorf("slug is too short, min length is 3")
	}

	if len(c.Slug) > 300 {
		return fmt.Errorf("slug is too long, max length is 300")
	}

	return nil
}

type SegmentCreated struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}
