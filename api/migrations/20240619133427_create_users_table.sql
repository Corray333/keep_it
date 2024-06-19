-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    user_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username text COLLATE pg_catalog."default" NOT NULL UNIQUE,
    email text COLLATE pg_catalog."default" NOT NULL UNIQUE,
    password character varying(60) COLLATE pg_catalog."default" NOT NULL,
    avatar text COLLATE pg_catalog."default",
    ref_code VARCHAR(6),
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT users_email_key UNIQUE (email)
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
