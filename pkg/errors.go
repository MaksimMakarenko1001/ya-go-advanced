package pkg

import (
	"fmt"
	"net/http"
)

var ErrInternalServer = &Error{
	Message: "Internal error",
	Code:    "INTERNAL_SERVER_ERROR",
	Status:  http.StatusInternalServerError,
}

var ErrNotFound = &Error{
	Message: "Not found",
	Code:    "NOT_FOUND",
	Status:  http.StatusNotFound,
}

var ErrBadRequest = &Error{
	Message: "Bad request",
	Code:    "BAD_REQUEST",
	Status:  http.StatusBadRequest,
}

var allowStatusError = map[int]struct{}{
	http.StatusInternalServerError: {},
	http.StatusNotFound:            {},
	http.StatusBadRequest:          {},
}

type ErrorCode string

type Error struct {
	Message string
	Code    ErrorCode
	Status  int
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
	if _, ok := allowStatusError[e.Status]; !ok {
		return http.StatusInternalServerError
	}
	return e.Status
}
