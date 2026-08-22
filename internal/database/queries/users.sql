-- name: CreateUser :execresult
INSERT INTO users (
    id,
    username,
    email,
    password_hash
)
VALUES (?, ?, ?, ?);

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;