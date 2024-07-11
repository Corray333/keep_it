-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    category_id INTEGER GENERATED ALWAYS AS IDENTITY (INCREMENT 1 START 1 MINVALUE 1 CACHE 1),
    owner BIGINT REFERENCES users(user_id),
    name VARCHAR(128),
    parent_category_id INTEGER,
    parent_owner BIGINT,
    FOREIGN KEY (parent_category_id, parent_owner) REFERENCES categories(category_id, owner),
    PRIMARY KEY (category_id, owner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd