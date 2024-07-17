-- +goose Up
-- +goose StatementBegin
ALTER TABLE note_tag DROP CONSTRAINT note_tag_tag_id_owner_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE note_tag ADD CONSTRAINT note_tag_tag_id_owner_fkey FOREIGN KEY (tag_id, owner) REFERENCES tag(id, owner) ON DELETE CASCADE;
-- +goose StatementEnd
