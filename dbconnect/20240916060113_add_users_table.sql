-- +goose Up
CREATE TABLE users (
    id bigserial PRIMARY KEY,
    Username text NOT NULL,
    Password text NOT NULL
);

-- +goose StatementBegin
SELECT 'Up migration applied: Table users created with columns id, Username, Password.';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS users;

-- +goose StatementBegin
SELECT 'Down migration applied: Table users dropped.';
-- +goose StatementEnd
