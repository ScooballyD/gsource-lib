// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: games.sql

package database

import (
	"context"
)

const addGame = `-- name: AddGame :one
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
RETURNING id, created_at, updated_at, title, url, image, category
`

type AddGameParams struct {
	Title    string
	Url      string
	Image    string
	Category string
}

func (q *Queries) AddGame(ctx context.Context, arg AddGameParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, addGame,
		arg.Title,
		arg.Url,
		arg.Image,
		arg.Category,
	)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Image,
		&i.Category,
	)
	return i, err
}

const getEpicGame = `-- name: GetEpicGame :one
SELECT title, url AS Href,
image, category
FROM games
WHERE category = 'Epic'
`

type GetEpicGameRow struct {
	Title    string
	Href     string
	Image    string
	Category string
}

func (q *Queries) GetEpicGame(ctx context.Context) (GetEpicGameRow, error) {
	row := q.db.QueryRowContext(ctx, getEpicGame)
	var i GetEpicGameRow
	err := row.Scan(
		&i.Title,
		&i.Href,
		&i.Image,
		&i.Category,
	)
	return i, err
}

const getGames = `-- name: GetGames :many
SELECT title, url AS Href,
image, category
FROM games
`

type GetGamesRow struct {
	Title    string
	Href     string
	Image    string
	Category string
}

func (q *Queries) GetGames(ctx context.Context) ([]GetGamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGamesRow
	for rows.Next() {
		var i GetGamesRow
		if err := rows.Scan(
			&i.Title,
			&i.Href,
			&i.Image,
			&i.Category,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const resetGames = `-- name: ResetGames :exec
DELETE FROM games
`

func (q *Queries) ResetGames(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetGames)
	return err
}
