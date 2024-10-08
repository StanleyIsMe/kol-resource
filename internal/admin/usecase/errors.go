package usecase

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

type DumplicatedUsernameError struct {
	username string
}

func (e DumplicatedUsernameError) ErrorCode() string {
	return "DUPLICATED_USERNAME"
}

func (e DumplicatedUsernameError) ErrorMsg() string {
	return fmt.Sprintf("username %s already exists", e.username)
}

func (e DumplicatedUsernameError) Error() string {
	return fmt.Sprintf("username %s already exists", e.username)
}

func (e DumplicatedUsernameError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type UnauthorizedError struct {
	err error
}

func (e UnauthorizedError) ErrorCode() string {
	return "UNAUTHORIZED"
}

func (e UnauthorizedError) ErrorMsg() string {
	return "Unauthorized"
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %v", e.err)
}

func (e UnauthorizedError) HTTPStatusCode() int {
	return http.StatusUnauthorized
}
