package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type UpdateUserSegments struct {
	JoinSegments  []string `json:"join"`
	LeaveSegments []string `json:"leave"`
}

func (u *UpdateUserSegments) CanJoin() bool {
	return len(u.JoinSegments) > 0
}

func (u *UpdateUserSegments) CanLeave() bool {
	return len(u.LeaveSegments) > 0
}

func (u *UpdateUserSegments) Validate() error {
	var errorBuilder strings.Builder
	if len(u.JoinSegments) == 0 && len(u.LeaveSegments) == 0 {
		errorBuilder.WriteString("no segments to join or leave\n")
	}

	var joinSegmentsMap = make(map[string]bool)
	for _, segment := range u.JoinSegments {
		if _, ok := joinSegmentsMap[segment]; ok {
			errorBuilder.WriteString(fmt.Sprintf("segment %s is duplicated in join segments\n", segment))
		}

		joinSegmentsMap[segment] = true
	}

	var leaveSegmentsMap = make(map[string]bool)
	for _, segment := range u.LeaveSegments {
		if _, ok := leaveSegmentsMap[segment]; ok {
			errorBuilder.WriteString(fmt.Sprintf("segment %s is duplicated in leave segments\n", segment))
		}

		leaveSegmentsMap[segment] = true
	}

	for _, segment := range u.JoinSegments {
		if _, ok := leaveSegmentsMap[segment]; ok {
			errorBuilder.WriteString(fmt.Sprintf("segment %s is in both join and leave segments\n", segment))
		}
	}

	if errorBuilder.Len() > 0 {
		return errors.New(errorBuilder.String())
	}

	return nil
}

type UserSegment struct {
	Slug string `json:"slug"`
}

type HistoryReport struct {
	UserID uint64 `json:"user_id"`
	Slug   string `json:"slug"`

	// OperationType is one of "join" or "leave".
	OperationType string `json:"operation_type"`

	// HappenedAt is a timestamp in RFC3339 format.
	HappenedAt string `json:"happened_at"`
}

type HistoryReportList struct {
	Segments []*HistoryReport
}

func (hrl *HistoryReportList) MarshalJSON() ([]byte, error) {
	return json.Marshal(hrl.Segments)
}

func (hrl *HistoryReportList) ToRawTable() [][]string {
	var table = make([][]string, len(hrl.Segments)+1)
	table[0] = []string{"user_id", "slug", "operation_type", "happened_at"}
	for i, segment := range hrl.Segments {
		table[i+1] = []string{
			fmt.Sprintf("%d", segment.UserID),
			segment.Slug,
			segment.OperationType,
			segment.HappenedAt,
		}
	}

	return table
}

type SegmentWithTTL struct {
	Slug string `json:"slug"`

	// TTL in days.
	TTL int `json:"ttl"`
}
