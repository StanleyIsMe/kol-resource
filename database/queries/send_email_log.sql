-- name: CreateSendEmailLog :one
INSERT INTO send_email_log (id, kol_id, email, admin_id, admin_name, product_id, product_name)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id, kol_id, email, admin_id, admin_name, product_id, product_name;
