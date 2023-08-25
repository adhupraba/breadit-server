-- name: FindVotesOfAPost :many
SELECT * FROM votes WHERE post_id = $1;

-- name: FindVotesOfPosts :many
SELECT * FROM votes WHERE post_id = ANY($1::INT[]);
