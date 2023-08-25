package dto

import (
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

type UserSegments struct {
	Segments []string `json:"segments"`
}
