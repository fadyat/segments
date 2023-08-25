package repository

import "avito-internship-2023/internal/api"

type Error interface {
	ToApiError() error
	error
}

type baseError struct {
	msg string
}

func (e *baseError) Error() string {
	return e.msg
}

type AlreadyExistsError struct {
	baseError
}

func NewAlreadyExistsError(msg string) *AlreadyExistsError {
	return &AlreadyExistsError{
		baseError: baseError{msg: msg},
	}
}

func (e *AlreadyExistsError) ToApiError() error {
	return api.NewConflictError(e.msg)
}

type NotFoundError struct {
	baseError
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		baseError: baseError{msg: msg},
	}
}

func (e *NotFoundError) ToApiError() error {
	return api.NewNotFoundError(e.msg)
}

type NotFoundMultiError struct {
	baseError
	args []string
}

func NewNotFoundMultiError(msg string, args ...string) *NotFoundMultiError {
	return &NotFoundMultiError{
		baseError: baseError{msg: msg},
		args:      args,
	}
}

func (e *NotFoundMultiError) ToApiError() error {
	return api.NewUnprocessableEntityError(e.msg, e.args...)
}
