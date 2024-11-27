-- name: AddDiscount :one
INSERT INTO discounts (id, created_at, updated_at, title, url, image, category, price, og_price, discount, rating)
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
    $7,
    $8
)
RETURNING *;

-- name: GetDiscounts :many
SELECT title, url AS Href,
image, category, price, og_price, discount, rating
FROM discounts;

-- name: GetDiscountsPrice :many
SELECT title, url AS Href,
image, category, price, og_price, discount, rating
FROM discounts
ORDER BY price;

-- name: GetDiscountsTitle :many
SELECT title, url AS Href,
image, category, price, og_price, discount, rating
FROM discounts
ORDER BY title;

-- name: GetDiscountsDiscount :many
SELECT title, url AS Href,
image, category, price, og_price, discount, rating
FROM discounts
ORDER BY discount DESC;

-- name: ResetDiscounts :exec
DELETE FROM discounts;