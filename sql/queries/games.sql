-- name: AddGame :one
INSERT INTO games (id, created_at, updated_at, title, url, image, category)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetGames :many
SELECT title, url AS Href,
image, category
FROM games;

-- name: GetEpicGame :one
SELECT title, url AS Href,
image, category
FROM games
WHERE category = 'Epic';

-- name: ResetGames :exec
DELETE FROM games;