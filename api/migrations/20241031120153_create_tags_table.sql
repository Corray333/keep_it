-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags(
    owner_id BIGINT REFERENCES users(user_id),
    tag_text VARCHAR(32) NOT NULL DEFAULT 'new tag',
    tag_color VARCHAR(8),
    PRIMARY KEY (owner_id, tag_text)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
