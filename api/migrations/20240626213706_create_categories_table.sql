-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    user_id BIGINT REF
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
