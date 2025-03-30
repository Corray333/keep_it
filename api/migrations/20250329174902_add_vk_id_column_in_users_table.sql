-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN vk_id BIGINT UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN vk_id;
-- +goose StatementEnd
