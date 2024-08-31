-- name: CreateProduct :one
INSERT INTO product (id, name, description, updated_admin_id)
  VALUES ($1, $2, $3, $4)
RETURNING
    id, name, description, updated_admin_id;          

-- name: GetProductByID :one
SELECT
    id, name, description, updated_admin_id
FROM
    product
WHERE
    id = $1;
    
-- name: ListProducts :many
SELECT
    id, name, description, updated_admin_id
FROM
    product
WHERE
    deleted_at IS NULL
ORDER BY
    id DESC;    
        

-- name: UpdateProduct :one
UPDATE product
SET
    name = $2, description = $3, updated_admin_id = $4
WHERE
    id = $1
RETURNING
    id, name, description, updated_admin_id;

-- name: UpdateProduct :one
