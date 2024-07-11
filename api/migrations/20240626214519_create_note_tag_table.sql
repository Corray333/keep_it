-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS note_tag (
    note_id UUID NOT NULL,
    tag_id BIGINT,
    owner BIGINT,
    FOREIGN KEY (tag_id, owner) REFERENCES tags(tag_id, owner),
    CONSTRAINT note_tag_pkey PRIMARY KEY (note_id, tag_id, owner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS note_tag;
-- +goose StatementEnd
