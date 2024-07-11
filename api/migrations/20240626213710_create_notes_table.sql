-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS notes(
    note_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator BIGINT NOT NULL REFERENCES users,
    title VARCHAR(256) NOT NULL DEFAULT '',
    source VARCHAR(16) NOT NULL DEFAULT 'keep_it',
    original JSONB NOT NULL DEFAULT '{"text":"Keep it"}'::JSONB,
    font VARCHAR(16) NOT NULL DEFAULT 'sans-serif',
    created_at BIGINT,
    copied_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    type SMALLINT NOT NULL DEFAULT 1,
    content TEXT NOT NULL DEFAULT '[]',
    cover TEXT NOT NULL DEFAULT '',
    checked BOOLEAN NOT NULL DEFAULT false,
    category_owner BIGINT,
    category_id INTEGER,
    FOREIGN KEY (category_id, category_owner) REFERENCES categories(category_id, owner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
