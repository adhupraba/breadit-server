-- name: FindUserVoteOfAComment :one
SELECT * FROM comment_votes WHERE comment_votes.user_id = $1 AND comment_votes.comment_id = $2;

-- name: CreateCommentVote :one
INSERT INTO comment_votes (comment_id, user_id, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCommentVote :exec
UPDATE comment_votes SET type = $1 WHERE id = $2;

-- name: RemoveCommentVote :exec
DELETE FROM comment_votes WHERE id = $1;
