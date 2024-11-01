-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    user_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT BY 1 START WITH 1 MINVALUE 1),
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    tg_username TEXT NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    ref_code VARCHAR(6) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
);
CREATE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE INDEX IF NOT EXISTS users_ref_code_idx ON users (ref_code);
CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP INDEX IF EXISTS users_username_idx;
DROP INDEX IF EXISTS users_ref_code_idx;
DROP INDEX IF EXISTS users_email_idx;
-- +goose StatementEnd
