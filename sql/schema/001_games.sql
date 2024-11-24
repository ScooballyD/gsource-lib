-- +goose Up
CREATE TABLE games(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    image TEXT NOT NULL,
    category TEXT NOT NULL
);

-- +goose Down
DROP TABLE games;