-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, username, password, image)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;