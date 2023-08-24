-- name: FindPostsOfASubreddit :many
SELECT
  posts.*,
  sqlc.embed(users),
  JSON_AGG(votes.*) AS votes,
  JSON_AGG(comments.*) AS comments
FROM posts
  INNER JOIN users ON users.id = posts.author_id
  LEFT JOIN votes ON votes.post_id = posts.id
  LEFT JOIN comments ON comments.post_id = posts.id
WHERE posts.subreddit_id = $1
GROUP BY posts.id, users.id
OFFSET $2 LIMIT $3;
