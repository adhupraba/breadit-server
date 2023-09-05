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
	UserID      int32 `db:"user_id" json:"userId"`
	SubredditID int32 `db:"subreddit_id" json:"subredditId"`
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

const findAllSubscriptionsOfUser = `-- name: FindAllSubscriptionsOfUser :many
SELECT id, user_id, subreddit_id, created_at, updated_at FROM subscriptions WHERE user_id = $1
`

func (q *Queries) FindAllSubscriptionsOfUser(ctx context.Context, userID int32) ([]Subscription, error) {
	rows, err := q.db.QueryContext(ctx, findAllSubscriptionsOfUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subscription
	for rows.Next() {
		var i Subscription
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SubredditID,
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

const findSubscriptionCountOfSubreddit = `-- name: FindSubscriptionCountOfSubreddit :one
SELECT COUNT(*) FROM subscriptions WHERE subreddit_id = $1
`

func (q *Queries) FindSubscriptionCountOfSubreddit(ctx context.Context, subredditID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, findSubscriptionCountOfSubreddit, subredditID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const findUserSubscription = `-- name: FindUserSubscription :one
SELECT id, user_id, subreddit_id, created_at, updated_at FROM subscriptions WHERE user_id = $1 AND subreddit_id = $2
`

type FindUserSubscriptionParams struct {
	UserID      int32 `db:"user_id" json:"userId"`
	SubredditID int32 `db:"subreddit_id" json:"subredditId"`
}

func (q *Queries) FindUserSubscription(ctx context.Context, arg FindUserSubscriptionParams) (Subscription, error) {
	row := q.db.QueryRowContext(ctx, findUserSubscription, arg.UserID, arg.SubredditID)
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

const removeSubscriptionUsingId = `-- name: RemoveSubscriptionUsingId :exec
DELETE FROM subscriptions WHERE id = $1
`

func (q *Queries) RemoveSubscriptionUsingId(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, removeSubscriptionUsingId, id)
	return err
}

const unsubscribeFromASubreddit = `-- name: UnsubscribeFromASubreddit :exec
DELETE FROM subscriptions WHERE user_id = $1 AND subreddit_id = $2
`

type UnsubscribeFromASubredditParams struct {
	UserID      int32 `db:"user_id" json:"userId"`
	SubredditID int32 `db:"subreddit_id" json:"subredditId"`
}

func (q *Queries) UnsubscribeFromASubreddit(ctx context.Context, arg UnsubscribeFromASubredditParams) error {
	_, err := q.db.ExecContext(ctx, unsubscribeFromASubreddit, arg.UserID, arg.SubredditID)
	return err
}
