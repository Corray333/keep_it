-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS notes(
    note_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner BIGINT NOT NULL REFERENCES users,
    title VARCHAR(256) NOT NULL DEFAULT '',
    source VARCHAR(8) NOT NULL DEFAULT 'keep_it',
    original JSONB NOT NULL DEFAULT '{"text":"Keep it"}'::JSONB,
    font VARCHAR(16) NOT NULL DEFAULT 'sans-serif',
    created_at BIGINT,
    copied_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    type SMALLINT NOT NULL DEFAULT 1,
    checked BOOLEAN NOT NULL DEFAULT false,
    categorie_id UUID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
