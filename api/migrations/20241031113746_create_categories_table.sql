-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    category_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    owner_id BIGINT REFERENCES users(user_id),
    name VARCHAR(128) NOT NULL,
    parent_category_id UUID,
    FOREIGN KEY (parent_category_id) REFERENCES categories(category_id),
    PRIMARY KEY (category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd