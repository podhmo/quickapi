package qdump

import (
	"errors"
	"fmt"
)

type StatusCoder interface {
	StatusCode() int
}

func StatusCodeOf(err error) int {
	code := 500
	var t StatusCoder
	if errors.As(err, &t) {
		code = t.StatusCode()
	}
	return code
}

func NewAPIError(err error, code int) interface {
	StatusCoder
	error
} {
	return &apiError{err: err, code: code}
}

type apiError struct {
	code int
	err  error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("api-error: %s", e.err.Error())
}

func (e *apiError) Unwrap() error {
	return e.err
}

func (e *apiError) StatusCode() int {
	return e.code
}
