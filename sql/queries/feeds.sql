-- name: CreateFeed :one
INSERT INTO feeds (ID, Created_At, Updated_At, Name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;