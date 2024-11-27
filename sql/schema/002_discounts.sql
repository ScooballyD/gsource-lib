-- +goose Up
CREATE TABLE discounts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    image TEXT NOT NULL,
    category TEXT NOT NULL,
    price FLOAT NOT NULL,
    og_price TEXT NOT NULL,
    discount TEXT NOT NULL,
    rating TEXT NOT NULL
);

-- +goose Down
DROP TABLE discounts;