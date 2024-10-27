-- +goose Up
-- +goose StatementBegin
CREATE TABLE "accounts" (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    balance DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

ALTER TABLE "accounts" ADD CONSTRAINT "unique_user_account" UNIQUE (user_id);

CREATE TABLE "transfers" (
    id BIGSERIAL PRIMARY KEY,
    from_account_id INTEGER NOT NULL,
    to_account_id INTEGER NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(from_account_id) REFERENCES accounts(id),
    FOREIGN KEY(to_account_id) REFERENCES accounts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "transfers";
ALTER TABLE "accounts" DROP CONSTRAINT "unique_user_account";
DROP TABLE IF EXISTS "accounts";
-- +goose StatementEnd
