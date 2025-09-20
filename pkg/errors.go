package pkg

import (
	"fmt"
	"net/http"
)

var ErrInternalServer = &Error{
	Message:    "Internal error",
	Code:       "INTERNAL_SERVER_ERROR",
	HttpStatus: http.StatusInternalServerError,
}

var ErrNotFound = &Error{
	Message:    "Not found",
	Code:       "NOT_FOUND",
	HttpStatus: http.StatusNotFound,
}

var ErrBadRequest = &Error{
	Message:    "Bad request",
	Code:       "BAD_REQUEST",
	HttpStatus: http.StatusBadRequest,
}

type ErrorCode string

type Error struct {
	Message    string
	Code       ErrorCode
	HttpStatus int
}

func (e *Error) Error() string {
	if e == nil || e.Code == "" {
		return ""
	}

	return fmt.Sprintf("[%s] %s", e.Code, e.Message)

}

func (e *Error) HTTPStatus() int {
	if e == nil {
		return http.StatusOK
	}
	if e.HttpStatus < 100 || e.HttpStatus >= 600 {
		return http.StatusInternalServerError
	}
	return e.HttpStatus
}
