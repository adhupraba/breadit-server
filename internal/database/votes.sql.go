// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: votes.sql

package database

import (
	"context"

	"github.com/lib/pq"
)

const createVote = `-- name: CreateVote :one
INSERT INTO votes (post_id, user_id, type)
VALUES ($1, $2, $3)
RETURNING id, post_id, user_id, type, created_at, updated_at
`

type CreateVoteParams struct {
	PostID int32    `db:"post_id" json:"postId"`
	UserID int32    `db:"user_id" json:"userId"`
	Type   VoteType `db:"type" json:"type"`
}

func (q *Queries) CreateVote(ctx context.Context, arg CreateVoteParams) (Vote, error) {
	row := q.db.QueryRowContext(ctx, createVote, arg.PostID, arg.UserID, arg.Type)
	var i Vote
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.UserID,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserVoteOfAPost = `-- name: FindUserVoteOfAPost :one
SELECT id, post_id, user_id, type, created_at, updated_at FROM votes WHERE post_id = $1 AND user_id = $2
`

type FindUserVoteOfAPostParams struct {
	PostID int32 `db:"post_id" json:"postId"`
	UserID int32 `db:"user_id" json:"userId"`
}

func (q *Queries) FindUserVoteOfAPost(ctx context.Context, arg FindUserVoteOfAPostParams) (Vote, error) {
	row := q.db.QueryRowContext(ctx, findUserVoteOfAPost, arg.PostID, arg.UserID)
	var i Vote
	err := row.Scan(
		&i.ID,
		&i.PostID,
		&i.UserID,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findVotesOfAPost = `-- name: FindVotesOfAPost :many
SELECT id, post_id, user_id, type, created_at, updated_at FROM votes WHERE post_id = $1
`

func (q *Queries) FindVotesOfAPost(ctx context.Context, postID int32) ([]Vote, error) {
	rows, err := q.db.QueryContext(ctx, findVotesOfAPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Vote
	for rows.Next() {
		var i Vote
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.UserID,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const findVotesOfPosts = `-- name: FindVotesOfPosts :many
SELECT id, post_id, user_id, type, created_at, updated_at FROM votes WHERE post_id = ANY($1::INT[])
`

func (q *Queries) FindVotesOfPosts(ctx context.Context, dollar_1 []int32) ([]Vote, error) {
	rows, err := q.db.QueryContext(ctx, findVotesOfPosts, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Vote
	for rows.Next() {
		var i Vote
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.UserID,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateVote = `-- name: UpdateVote :exec
UPDATE votes SET type = $1 WHERE id = $2
`

type UpdateVoteParams struct {
	Type VoteType `db:"type" json:"type"`
	ID   int32    `db:"id" json:"id"`
}

func (q *Queries) UpdateVote(ctx context.Context, arg UpdateVoteParams) error {
	_, err := q.db.ExecContext(ctx, updateVote, arg.Type, arg.ID)
	return err
}
