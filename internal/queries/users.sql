-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: FindUserById :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, username, password, image)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUsername :exec
UPDATE users SET username = $1 WHERE id = $2;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username = $1;
