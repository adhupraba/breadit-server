-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, subreddit_id) VALUES ($1, $2) RETURNING *;

-- name: FindUserSubscription :one
SELECT * FROM subscriptions WHERE user_id = $1 AND subreddit_id = $2;

-- name: FindSubscriptionCountOfSubreddit :one
SELECT COUNT(*) FROM subscriptions WHERE subreddit_id = $1;

-- name: RemoveSubscriptionUsingId :exec
DELETE FROM subscriptions WHERE id = $1;

-- name: UnsubscribeFromASubreddit :exec
DELETE FROM subscriptions WHERE user_id = $1 AND subreddit_id = $2;

-- name: FindAllSubscriptionsOfUser :many
SELECT * FROM subscriptions WHERE user_id = $1;
