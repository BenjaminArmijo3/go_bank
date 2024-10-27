-- +goose Up
-- +goose StatementBegin

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "email" varchar(256) UNIQUE NOT NULL,
  "hashed_password" varchar(256) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
