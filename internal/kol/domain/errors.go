package domain

import (
	"errors"
	"fmt"
)

var ErrDataNotFound = errors.New("data not found")

type GenerateUUIDError struct {
	Err error
}

func (e GenerateUUIDError) Error() string {
	return fmt.Sprintf("failed to generate uuid: %s", e.Err)
}

type InsertRecordError struct {
	Err error
}

func (e InsertRecordError) Error() string {
	return fmt.Sprintf("failed to insert record into db: %s", e.Err)
}

type QueryRecordError struct {
	Err error
}

func (e QueryRecordError) Error() string {
	return fmt.Sprintf("failed to query db: %s", e.Err)
}

type UpdateRecordError struct {
	Err error
}

func (e UpdateRecordError) Error() string {
	return fmt.Sprintf("failed to update record in db: %s", e.Err)
}

type DeleteRecordError struct {
	Err error
}

func (e DeleteRecordError) Error() string {
	return fmt.Sprintf("failed to delete record in db: %s", e.Err)
}

type UUIDInvalidError struct {
	Field string
	UUID  string
}

func (m UUIDInvalidError) Error() string {
	return fmt.Sprintf("%s is %s not a valid uuid", m.Field, m.UUID)
}

type SexInvalidError struct {
	Sex string
}

func (m SexInvalidError) Error() string {
	return fmt.Sprintf("%s is not a valid sex type", m.Sex)
}
