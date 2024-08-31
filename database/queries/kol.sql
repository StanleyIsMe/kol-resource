-- name: CreateKol :one
INSERT INTO kol (id, name, email, description, sex, enable, updated_admin_id)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id, name, email, description, sex, enable, updated_admin_id;    

-- name: GetKolByEmail :one
SELECT
    id, name, email, description, sex, enable, updated_admin_id
FROM
    kol
WHERE
    email = $1;
    
-- name: GetKolByID :one
SELECT
    id, name, email, description, sex, enable, updated_admin_id
FROM
    kol
WHERE
    id = $1;
    
-- name: ListKols :many
SELECT
    id, name, email, description, sex, enable, updated_admin_id
FROM
    kol
WHERE
    deleted_at IS NULL
ORDER BY
    id DESC;
    

-- name: UpdateKol :one
UPDATE kol
SET
    name = $2, email = $3, description = $4, sex = $5, enable = $6, updated_admin_id = $7
WHERE
    id = $1
RETURNING
    id, name, email, description, sex, enable, updated_admin_id;

-- name: UpdateKol :one
