-- name: FindPostsOfASubreddit :many
SELECT * FROM posts WHERE subreddit_id = $1 OFFSET $2 LIMIT $3;

-- name: PostsData :many
SELECT
  posts.*,
  sqlc.embed(users),
  json_agg(votes.*) AS votes,
  json_agg(comments.*) AS comments
FROM posts
  INNER JOIN users ON users.id = posts.author_id
  LEFT JOIN votes ON votes.post_id = posts.id
  LEFT JOIN comments ON comments.post_id = posts.id
WHERE posts.subreddit_id = $1
GROUP BY posts.id, users.id;