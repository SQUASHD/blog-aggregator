-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES ($1, $2, $3, $4, $5, $6, null)
RETURNING *;
-- name: GetAllFeeds :many
SELECT * FROM feeds;
-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;
-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = $1 WHERE id = $2;