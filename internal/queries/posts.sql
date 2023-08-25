-- name: FindPostsOfSubredditWithAuthor :many
SELECT
  sqlc.embed(posts),
  sqlc.embed(users)
FROM posts
  INNER JOIN users ON users.id = posts.author_id
WHERE posts.subreddit_id = $1
GROUP BY posts.id, users.id
OFFSET $2 LIMIT $3;
