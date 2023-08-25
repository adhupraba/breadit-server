-- name: CreateSubreddit :one
INSERT INTO subreddits (name, creator_id) VALUES ($1, $2) RETURNING *;

-- name: FindSubredditByName :one
SELECT * FROM subreddits WHERE name = $1;

-- name: FindSubredditOfCreator :one
SELECT * FROM subreddits WHERE id = $1 AND creator_id = $2;
