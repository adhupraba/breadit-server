-- name: CreateComment :one
INSERT INTO comments (text, post_id, author_id, reply_to_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FindCommentsOfAPost :many
SELECT
  sqlc.embed(comments),
  sqlc.embed(users),
  TO_JSON(ARRAY_AGG(DISTINCT comment_votes.*)) AS votes
FROM comments
  INNER JOIN users ON users.id = comments.author_id
  LEFT JOIN comment_votes ON comment_votes.comment_id = comments.id
WHERE
  comments.post_id = $1 AND
  comments.reply_to_id IS NULL
GROUP BY comments.id, users.id;

-- name: FindRepliesForComments :many
SELECT
  sqlc.embed(comments),
  sqlc.embed(users),
  TO_JSON(ARRAY_AGG(DISTINCT comment_votes.*)) AS votes
FROM comments
  INNER JOIN users ON users.id = comments.author_id
  LEFT JOIN comment_votes ON comment_votes.comment_id = comments.id
WHERE
  comments.reply_to_id = ANY(sqlc.arg(comment_ids)::INT[])
GROUP BY comments.id, users.id;