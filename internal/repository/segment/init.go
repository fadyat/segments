package segment

import (
	"avito-internship-2023/internal/repository"
	"errors"
	"github.com/lib/pq"
)

type repo struct {
	repository.Transactor
}

func NewRepository(
	transactor repository.Transactor,
) repository.ISegment {
	return &repo{
		Transactor: transactor,
	}
}

func toSegmentError(err error) error {
	var e *pq.Error
	switch {
	case errors.As(err, &e):
		return postgresErrorToCustomError(err.(*pq.Error))
	default:
		return err
	}
}

func postgresErrorToCustomError(err *pq.Error) error {
	switch err.Code {
	case "23505":
		return repository.NewAlreadyExistsError("segment already exists")
	default:
		return err
	}
}
