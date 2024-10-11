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

type DuplicatedResourceError struct {
	name     string
	resource string
}

func (e DuplicatedResourceError) ErrorCode() string {
	return "DUPLICATED_RESOURCE"
}

func (e DuplicatedResourceError) ErrorMsg() string {
	return fmt.Sprintf("%s %s already exists", e.resource, e.name)
}

func (e DuplicatedResourceError) Error() string {
	return fmt.Sprintf("%s %s already exists", e.resource, e.name)
}

func (e DuplicatedResourceError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type NotFoundError struct {
	resource string
	id       interface{}
}

func (e NotFoundError) ErrorCode() string {
	return "NOT_FOUND"
}

func (e NotFoundError) ErrorMsg() string {
	return fmt.Sprintf("%s %v not found", e.resource, e.id)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s %v not found", e.resource, e.id)
}

func (e NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}
