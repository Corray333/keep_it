-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_token (
    user_id INTEGER NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    token VARCHAR(512) NOT NULL,
    expires_at BIGINT NOT NULL,
    CONSTRAINT user_token_pk PRIMARY KEY (user_id, token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_token;
-- +goose StatementEnd
