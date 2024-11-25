-- name: AddDiscount :one
INSERT INTO discounts (id, created_at, updated_at, title, url, image, category, price, og_price, discount)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetDiscounts :many
SELECT title, url AS Href,
image, category, price, og_price, discount
FROM discounts;

-- name: ResetDiscounts :exec
DELETE FROM discounts;