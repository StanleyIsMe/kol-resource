// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlcdb

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Sex string

const (
	SexM Sex = "m"
	SexF Sex = "f"
)

func (e *Sex) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Sex(s)
	case string:
		*e = Sex(s)
	default:
		return fmt.Errorf("unsupported scan type for Sex: %T", src)
	}
	return nil
}

type NullSex struct {
	Sex   Sex
	Valid bool // Valid is true if Sex is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullSex) Scan(value interface{}) error {
	if value == nil {
		ns.Sex, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Sex.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullSex) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Sex), nil
}

func (e Sex) Valid() bool {
	switch e {
	case SexM,
		SexF:
		return true
	}
	return false
}

func AllSexValues() []Sex {
	return []Sex{
		SexM,
		SexF,
	}
}

type Admin struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Password  string
	Salt      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type Kol struct {
	ID             uuid.UUID
	Name           string
	Email          string
	Description    string
	Sex            Sex
	Enable         bool
	UpdatedAdminID uuid.UUID
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
	DeletedAt      pgtype.Timestamp
}

type KolTag struct {
	ID             uuid.UUID
	KolID          uuid.UUID
	TagID          uuid.UUID
	UpdatedAdminID uuid.UUID
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
	DeletedAt      pgtype.Timestamp
}

type Product struct {
	ID             uuid.UUID
	Name           string
	Description    string
	UpdatedAdminID uuid.UUID
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
	DeletedAt      pgtype.Timestamp
}

type SendEmailLog struct {
	ID          uuid.UUID
	KolID       uuid.UUID
	Email       uuid.UUID
	AdminID     uuid.UUID
	AdminName   string
	ProductID   uuid.UUID
	ProductName string
	CreatedAt   pgtype.Timestamp
}

type Tag struct {
	ID             uuid.UUID
	Name           string
	UpdatedAdminID uuid.UUID
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
	DeletedAt      pgtype.Timestamp
}
