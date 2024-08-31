-- name: CreateTag :one
INSERT INTO tag (id, name, updated_admin_id)
  VALUES ($1, $2, $3)
RETURNING
    id, name, updated_admin_id;

-- name: GetTagByID :one
SELECT
    id, name, updated_admin_id
FROM
    tag
WHERE
    id = $1;

-- name: ListTags :many
SELECT
    id, name, updated_admin_id
FROM
    tag
WHERE
    deleted_at IS NULL
ORDER BY
    id DESC;

-- name: UpdateTag :one
UPDATE tag
SET 
    name = $2, updated_admin_id = $3
WHERE
    id = $1
RETURNING
    id, name, updated_admin_id;
