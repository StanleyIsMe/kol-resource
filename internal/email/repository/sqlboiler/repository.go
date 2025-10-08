package sqlboiler

import (
	"context"
	"database/sql"
	"kolresource/internal/email/domain"
)

type EmailRepository struct {
	db *sql.DB
}

var _ domain.Repository = (*EmailRepository)(nil)

func NewEmailRepository(db *sql.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

