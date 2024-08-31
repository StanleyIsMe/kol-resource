-- name: CreateAdmin :one
INSERT INTO admin (id, name, username, password, salt)
  VALUES ($1, $2, $3, $4, $5)
RETURNING
    id, name;

-- name: GetAdminByUsername :one
SELECT
    id, name, username, password, salt
FROM
    admin
WHERE
    username = $1;

