-- name: FindPostsOfSubreddit :many
WITH vars (subreddit_name, is_authenticated, subreddit_ids, subreddit_id) AS (
	VALUES (
    sqlc.arg(subreddit_name)::TEXT,
    sqlc.arg(is_authenticated)::BOOL,
    sqlc.slice(subreddit_ids)::INT[],
    sqlc.arg(subreddit_id)::INT
  )
)
SELECT
  sqlc.embed(posts),
  sqlc.embed(users),
  sqlc.embed(subreddits),
  TO_JSON(ARRAY_AGG(DISTINCT votes.*)) AS votes,
  TO_JSON(ARRAY_AGG(DISTINCT comments.*)) AS comments
FROM posts
  INNER JOIN users ON users.id = posts.author_id
  INNER JOIN subreddits ON subreddits.id = posts.subreddit_id
  LEFT JOIN votes ON votes.post_id = posts.id
  -- take top level comments only
  LEFT JOIN comments ON comments.post_id = posts.id AND comments.reply_to_id IS NULL,
  vars
WHERE (
  CASE
    WHEN
      vars.subreddit_name IS NOT NULL AND
      LENGTH(vars.subreddit_name) > 0
        THEN subreddits.name = vars.subreddit_name
    WHEN 
      vars.is_authenticated AND
      ARRAY_LENGTH(vars.subreddit_ids, 1) > 0
        THEN subreddits.id = ANY(vars.subreddit_ids)
    WHEN
      vars.subreddit_id IS NOT NULL AND
      vars.subreddit_id > 0
        THEN posts.subreddit_id = vars.subreddit_id
 		ELSE TRUE
  END
)
GROUP BY posts.id, users.id, subreddits.id
ORDER BY posts.created_at DESC, posts.id DESC
OFFSET sqlc.arg('offset') LIMIT sqlc.arg('limit');

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