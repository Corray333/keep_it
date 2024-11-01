-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags(
    owner_id BIGINT REFERENCES users(user_id),
    text VARCHAR(32) NOT NULL DEFAULT 'new tag',
    color VARCHAR(6),
    PRIMARY KEY (owner_id, text)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
