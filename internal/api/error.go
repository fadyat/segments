package api

import (
	"net/http"
	"strings"
)

type Error interface {
	StatusCode() int
	error
}

type baseError struct {
	Msg string `json:"error"`
}

func (e *baseError) Error() string {
	return e.Msg
}

type BadRequestError struct {
	baseError
}

func NewBadRequestError(msg string) *BadRequestError {
	return &BadRequestError{
		baseError: baseError{Msg: msg},
	}
}

func (e *BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

type NotFoundError struct {
	baseError
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		baseError: baseError{Msg: msg},
	}
}

func (e *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

type ConflictError struct {
	baseError
}

func NewConflictError(msg string) *ConflictError {
	return &ConflictError{
		baseError: baseError{Msg: msg},
	}
}

func (e *ConflictError) StatusCode() int {
	return http.StatusConflict
}

type UnprocessableEntityError struct {
	baseError
	Fields []string `json:"fields"`
}

func (e *UnprocessableEntityError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

func takeNonEmpty(fields []string) []string {
	var nonEmpty = make([]string, 0, len(fields))
	for _, field := range fields {
		if field != "" {
			nonEmpty = append(nonEmpty, field)
		}
	}

	return nonEmpty
}

func NewUnprocessableEntityError(msg string, fields ...string) *UnprocessableEntityError {
	var byNewLine = make([]string, 0, len(fields))
	for _, field := range fields {
		byNewLine = append(
			byNewLine,
			takeNonEmpty(strings.Split(field, "\n"))...,
		)
	}

	return &UnprocessableEntityError{
		baseError: baseError{Msg: msg},
		Fields:    byNewLine,
	}
}
