-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetFeeds :many
SELECT * 
FROM feeds;

-- name: GetFeedsWithUsernames :many
SELECT f.id, f.created_at, f.updated_at, f.name, f.url, f.user_id, u.name as user_name
FROM feeds f
JOIN users u ON f.user_id = u.id;

-- name: CreateFollowFeed :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    ins.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow ins
INNER JOIN users ON ins.user_id = users.id
INNER JOIN feeds ON ins.feed_id = feeds.id;

-- name: GetFeedByURL :one
SELECT id, name
FROM feeds
WHERE url = $1;