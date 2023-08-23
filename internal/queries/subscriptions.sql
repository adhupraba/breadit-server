-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, subreddit_id) VALUES ($1, $2) RETURNING *;