-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notes(
    note_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    title VARCHAR(256) NOT NULL DEFAULT '',
    source VARCHAR(16) NOT NULL DEFAULT 'keep_it',
    original TEXT NOT NULL DEFAULT '',
    icon JSONB NOT NULL DEFAULT '{"type":"custom","icon":"keep-it"}'::JSONB,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    copied_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    type SMALLINT NOT NULL DEFAULT 1,
    content JSON NOT NULL DEFAULT '[]',
    cover TEXT NOT NULL DEFAULT '',
    checked BOOLEAN NOT NULL DEFAULT false,
    category_id UUID,
    FOREIGN KEY (category_id) REFERENCES categories(category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
-- +goose StatementEnd
