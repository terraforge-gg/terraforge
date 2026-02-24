-- +goose Up
-- +goose StatementBegin
CREATE TABLE "jwks" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "publicKey" TEXT NOT NULL,
    "privateKey" TEXT NOT NULL,
    "createdAt" timestamptz NOT NULL,
    "expiresAt" timestamptz
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE "jwks";
-- +goose StatementEnd