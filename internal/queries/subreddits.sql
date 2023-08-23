-- name: FindSubredditByName :one
SELECT * FROM subreddits WHERE name = $1;

-- name: CreateSubreddit :one
INSERT INTO subreddits (name, creator_id) VALUES ($1, $2) RETURNING *;
