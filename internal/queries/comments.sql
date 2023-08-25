-- name: FindCommentsOfAPost :many
SELECT * FROM comments WHERE post_id = $1;

-- name: FindCommentsOfPosts :many
SELECT * FROM comments WHERE post_id = ANY($1::INT[]);
