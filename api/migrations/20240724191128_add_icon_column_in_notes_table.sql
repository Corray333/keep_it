-- +goose Up
-- +goose StatementBegin
ALTER TABLE notes ADD COLUMN icon JSONB NOT NULL DEFAULT '{"type":"custom","icon":"keep-it"}'::JSONB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notes DROP COLUMN icon;
-- +goose StatementEnd
