package shared

import (
	"errors"
	"fmt"
)

var (
	API_ERROR_FORMAT = "api-error: %s"
)

type StatusCoder interface {
	StatusCode() int
}

func StatusCodeOf(err error) int {
	return StatusCodeOfOrDefault(err, 500)
}
func StatusCodeOfOrDefault(err error, code int) int {
	var t StatusCoder
	if errors.As(err, &t) {
		code = t.StatusCode()
	}
	return code
}

type APIError struct {
	code int
	err  error
}

func NewAPIError(err error, code int) *APIError {
	return &APIError{err: err, code: code}
}

func (e *APIError) Error() string {
	return fmt.Sprintf(API_ERROR_FORMAT, e.err.Error())
}

func (e *APIError) Unwrap() error {
	return e.err
}

func (e *APIError) StatusCode() int {
	return e.code
}

var _ interface {
	StatusCoder
	error
} = &APIError{}
