-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    category_id INTEGER GENERATED ALWAYS AS IDENTITY (INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    user_id BIGINT REFERENCES users(user_id),
    name VARCHAR(128),
    parent_category_id INTEGER,
    parent_user_id BIGINT,
    FOREIGN KEY (parent_category_id, parent_user_id) REFERENCES categories(category_id, user_id),
    PRIMARY KEY (category_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd