-- +goose Up
CREATE TABLE items (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    amount int NOT NULL,
    quantity int NOT NULL,
    status text NOT NULL,
    owner_id bigint NOT NULL
);

-- +goose StatementBegin
SELECT 'Up migration applied: Table items created with columns id, title, amount, quantity, status, and owner_id.';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS items;

-- +goose StatementBegin
SELECT 'Down migration applied: Table items dropped.';
-- +goose StatementEnd
