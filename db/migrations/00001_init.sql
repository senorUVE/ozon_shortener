-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url (
    id BIGSERIAL PRIMARY KEY,
    original_url varchar(2048) NOT NULL UNIQUE,
    token VARCHAR(10) UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url CASCADE;
-- +goose StatementEnd