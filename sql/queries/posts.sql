-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsFromUser :many
SELECT COUNT(p.*), p.*
FROM posts AS p
INNER JOIN
feeds AS f ON p.feed_id = f.id
WHERE f.user_id = $1
ORDER BY p.published_at DESC;