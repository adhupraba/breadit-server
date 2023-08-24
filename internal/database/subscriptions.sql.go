// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: subscriptions.sql

package database

import (
	"context"
)

const createSubscription = `-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, subreddit_id) VALUES ($1, $2) RETURNING id, user_id, subreddit_id, created_at, updated_at
`

type CreateSubscriptionParams struct {
	UserID      int32 `json:"userId"`
	SubredditID int32 `json:"subredditId"`
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, createSubscription, arg.UserID, arg.SubredditID)
	var i Subscription
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SubredditID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
