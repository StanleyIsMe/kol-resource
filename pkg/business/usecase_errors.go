package business

import (
	"fmt"
	"net/http"
)

type UseCaseError interface {
	ErrorCode() string
	ErrorMsg() string
	Error() string
	HTTPStatusCode() int
}

type InternalServerError struct {
	err error
}

func (e InternalServerError) ErrorCode() string {
	return "INTERNAL_SERVER_ERROR"
}

func (e InternalServerError) ErrorMsg() string {
	return "Internal Server Error"
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %v", e.err)
}

func (e InternalServerError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}
