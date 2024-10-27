-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTransferByID :one
SELECT * FROM transfers WHERE id = $1;

-- name: GetTransfersByFromAccountID :many
SELECT * FROM transfers where from_account_id = $1;

-- name: GetTransfersByToAccountID :many
SELECT * FROM transfers where from_account_id = $1;
