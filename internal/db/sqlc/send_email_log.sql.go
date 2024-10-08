// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: send_email_log.sql

package sqlcdb

import (
	"context"

	"github.com/google/uuid"
)

const CreateSendEmailLog = `-- name: CreateSendEmailLog :one
INSERT INTO send_email_log (id, kol_id, email, admin_id, admin_name, product_id, product_name)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id, kol_id, email, admin_id, admin_name, product_id, product_name
`

type CreateSendEmailLogParams struct {
	ID          uuid.UUID
	KolID       uuid.UUID
	Email       uuid.UUID
	AdminID     uuid.UUID
	AdminName   string
	ProductID   uuid.UUID
	ProductName string
}

type CreateSendEmailLogRow struct {
	ID          uuid.UUID
	KolID       uuid.UUID
	Email       uuid.UUID
	AdminID     uuid.UUID
	AdminName   string
	ProductID   uuid.UUID
	ProductName string
}

func (q *Queries) CreateSendEmailLog(ctx context.Context, arg *CreateSendEmailLogParams) (*CreateSendEmailLogRow, error) {
	row := q.db.QueryRow(ctx, CreateSendEmailLog,
		arg.ID,
		arg.KolID,
		arg.Email,
		arg.AdminID,
		arg.AdminName,
		arg.ProductID,
		arg.ProductName,
	)
	var i CreateSendEmailLogRow
	err := row.Scan(
		&i.ID,
		&i.KolID,
		&i.Email,
		&i.AdminID,
		&i.AdminName,
		&i.ProductID,
		&i.ProductName,
	)
	return &i, err
}
