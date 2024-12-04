-- +goose Up
-- +goose StatementBegin
ALTER TABLE notes
    DROP CONSTRAINT notes_category_id_fkey,
    ADD CONSTRAINT notes_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notes
    DROP CONSTRAINT notes_category_id_fkey,
    ADD CONSTRAINT notes_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories(category_id);
-- +goose StatementEnd