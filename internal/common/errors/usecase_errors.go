package errors

import (
	"fmt"
	"net/http"
)

type DuplicatedResourceError struct {
	Name     string
	Resource string
}

func (e DuplicatedResourceError) ErrorCode() string {
	return "DUPLICATED_RESOURCE"
}

func (e DuplicatedResourceError) ErrorMsg() string {
	return fmt.Sprintf("%s %s already exists", e.Resource, e.Name)
}

func (e DuplicatedResourceError) Error() string {
	return fmt.Sprintf("%s %s already exists", e.Resource, e.Name)
}

func (e DuplicatedResourceError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e NotFoundError) ErrorCode() string {
	return "NOT_FOUND"
}

func (e NotFoundError) ErrorMsg() string {
	return fmt.Sprintf("%s %v not found", e.Resource, e.ID)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s %v not found", e.Resource, e.ID)
}

func (e NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}
