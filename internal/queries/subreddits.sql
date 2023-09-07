-- name: CreateSubreddit :one
INSERT INTO subreddits (name, creator_id) VALUES ($1, $2) RETURNING *;

-- name: FindSubredditByName :one
SELECT * FROM subreddits WHERE name = $1;

-- name: FindSubredditOfCreator :one
SELECT * FROM subreddits WHERE id = $1 AND creator_id = $2;

-- name: SearchSubreddits :many
SELECT
  sqlc.embed(subre),
  COALESCE(post_agg.post_count, 0) AS post_count,
  COALESCE(sub_agg.sub_count, 0) AS sub_count
FROM subreddits AS subre
LEFT JOIN (
  SELECT subreddit_id, COUNT(*) AS post_count
  FROM posts
  GROUP BY subreddit_id
) AS post_agg ON post_agg.subreddit_id = subre.id
LEFT JOIN (
  SELECT subreddit_id, COUNT(*) AS sub_count
  FROM subscriptions
  GROUP BY subreddit_id
) AS sub_agg ON sub_agg.subreddit_id = subre.id
WHERE subre.name LIKE $1
ORDER BY subre.name ASC
LIMIT 5;
