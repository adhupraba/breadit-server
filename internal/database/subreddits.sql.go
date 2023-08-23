// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: subreddits.sql

package database

import (
	"context"
	"database/sql"
)

const createSubreddit = `-- name: CreateSubreddit :one
INSERT INTO subreddits (name, creator_id) VALUES ($1, $2) RETURNING id, name, creator_id, created_at, updated_at
`

type CreateSubredditParams struct {
	Name      string
	CreatorID sql.NullInt32
}

func (q *Queries) CreateSubreddit(ctx context.Context, arg CreateSubredditParams) (Subreddit, error) {
	row := q.db.QueryRowContext(ctx, createSubreddit, arg.Name, arg.CreatorID)
	var i Subreddit
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatorID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findSubredditByName = `-- name: FindSubredditByName :one
SELECT id, name, creator_id, created_at, updated_at FROM subreddits WHERE name = $1
`

func (q *Queries) FindSubredditByName(ctx context.Context, name string) (Subreddit, error) {
	row := q.db.QueryRowContext(ctx, findSubredditByName, name)
	var i Subreddit
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatorID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
