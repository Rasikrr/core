package http

import (
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) StatusCode() int {
	return e.Code
}

func (e *Error) Wrap(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %v", e.Message, err),
	}
}

func NewError(message string, code int) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
