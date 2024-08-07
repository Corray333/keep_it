-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_note_access (
    user_id BIGINT NOT NULL,
    note_id UUID NOT NULL,
    PRIMARY KEY (user_id, note_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (note_id) REFERENCES notes (note_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_note_access;
-- +goose StatementEnd
