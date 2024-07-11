-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags(
    tag_id INTEGER GENERATED ALWAYS AS IDENTITY (INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    owner BIGINT REFERENCES users(user_id),
    text VARCHAR(16) NOT NULL DEFAULT 'new tag',
    color VARCHAR(6),
    PRIMARY KEY (tag_id, owner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
