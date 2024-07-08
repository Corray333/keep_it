-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags(
    tag_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    text VARCHAR(16) NOT NULL DEFAULT 'new tag',
    color VARCHAR(6)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
