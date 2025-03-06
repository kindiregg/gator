-- name: CreateUser :one
INSERT INTO users (ID, Created_At, Updated_At, Name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE Name = $1;

-- name: ResetUsers :exec
TRUNCATE TABLE users;

-- name: GetUsers :many
SELECT Name FROM users;