-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS note_tag (
    note_id UUID NOT NULL,
    tag_text VARCHAR(32) NOT NULL,
    owner_id BIGINT NOT NULL,
    FOREIGN KEY (tag_text, owner_id) REFERENCES tags(tag_text, owner_id) ON DELETE CASCADE,
    FOREIGN KEY (note_id) REFERENCES notes(note_id) ON DELETE CASCADE,
    CONSTRAINT note_tag_pkey PRIMARY KEY (note_id, tag_text, owner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS note_tag;
-- +goose StatementEnd
