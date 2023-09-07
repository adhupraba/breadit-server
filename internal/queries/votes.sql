-- name: FindUserVoteOfAPost :one
SELECT * FROM votes WHERE post_id = $1 AND user_id = $2;

-- name: CreatePostVote :one
INSERT INTO votes (post_id, user_id, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdatePostVote :exec
UPDATE votes SET type = $1 WHERE id = $2;

-- name: RemovePostVote :exec
DELETE FROM votes WHERE id = $1;
