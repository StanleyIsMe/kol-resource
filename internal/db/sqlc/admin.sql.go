// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: admin.sql

package sqlcdb

import (
	"context"

	"github.com/google/uuid"
)

const CreateAdmin = `-- name: CreateAdmin :one
INSERT INTO admin (id, name, username, password, salt)
  VALUES ($1, $2, $3, $4, $5)
RETURNING
    id, name
`

type CreateAdminParams struct {
	ID       uuid.UUID
	Name     string
	Username string
	Password string
	Salt     string
}

type CreateAdminRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) CreateAdmin(ctx context.Context, arg *CreateAdminParams) (*CreateAdminRow, error) {
	row := q.db.QueryRow(ctx, CreateAdmin,
		arg.ID,
		arg.Name,
		arg.Username,
		arg.Password,
		arg.Salt,
	)
	var i CreateAdminRow
	err := row.Scan(&i.ID, &i.Name)
	return &i, err
}

const GetAdminByUsername = `-- name: GetAdminByUsername :one
SELECT
    id, name, username, password, salt
FROM
    admin
WHERE
    username = $1
`

type GetAdminByUsernameRow struct {
	ID       uuid.UUID
	Name     string
	Username string
	Password string
	Salt     string
}

func (q *Queries) GetAdminByUsername(ctx context.Context, username string) (*GetAdminByUsernameRow, error) {
	row := q.db.QueryRow(ctx, GetAdminByUsername, username)
	var i GetAdminByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Password,
		&i.Salt,
	)
	return &i, err
}