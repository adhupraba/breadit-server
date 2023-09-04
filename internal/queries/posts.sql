-- name: FindPostsOfSubredditWithAuthor :many
SELECT
  sqlc.embed(posts),
  sqlc.embed(users)
FROM posts
  INNER JOIN users ON users.id = posts.author_id
WHERE posts.subreddit_id = $1
GROUP BY posts.id, users.id
OFFSET $2 LIMIT $3;

-- name: CreatePost :one
INSERT INTO posts (title, content, subreddit_id, author_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindPostWithAuthorAndVotes :one
SELECT
  sqlc.embed(posts),
  sqlc.embed(users),
  JSON_AGG(votes.*) AS votes
FROM posts
  INNER JOIN users ON users.id = posts.author_id
  LEFT JOIN votes ON votes.post_id = posts.id
WHERE posts.id = $1
GROUP BY posts.id, users.id;
