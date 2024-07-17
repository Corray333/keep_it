-- +goose Up
-- +goose StatementBegin
CREATE VIEW note_search AS
    SELECT notes.note_id, creator, title, source, original, font, created_at, copied_at, type, content, cover, checked, category_owner, category_id, tag_id, text, color
    FROM user_note_access 
    NATURAL JOIN notes 
    LEFT JOIN (note_tag NATURAL JOIN tags) AS nt ON nt.note_id = notes.note_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS note_search;
-- +goose StatementEnd
