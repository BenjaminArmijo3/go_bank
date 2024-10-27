-- name: CreateAccount :one
INSERT INTO accounts (
    user_id,
    balance
) VALUES ($1,$2) RETURNING *;


-- name: GetAccountById :one
SELECT * FROM accounts WHERE id= $1;


-- name: GetAccountByUserId :one
SELECT * FROM accounts WHERE user_id= $1;


-- name: UpdateAccountBalance :one
UPDATE accounts SET balance = $1 WHERE id= $2 RETURNING *;


-- name: UpdateAccountBalanceNew :one
UPDATE accounts SET balance = balance + sqlc.arg(amount) WHERE id = sqlc.arg(id) RETURNING *;
