-- name: FindVotesOfAPost :many
SELECT * FROM votes WHERE post_id = $1;

-- name: FindVotesOfPosts :many
SELECT * FROM votes WHERE post_id = ANY($1::INT[]);

-- name: FindUserVoteOfAPost :one
SELECT * FROM votes WHERE post_id = $1 AND user_id = $2;

-- name: CreateVote :one
INSERT INTO votes (post_id, user_id, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateVote :exec
UPDATE votes SET type = $1 WHERE id = $2;

-- name: RemoveVote :exec
DELETE FROM votes WHERE id = $1;
