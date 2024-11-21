-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users
    (name, password)
VALUES
    ($1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
    set name = $2,
    password = $3,
    updated_at = current_timestamp
WHERE
    id = $1
RETURNING *;

-- name: GetEventsFrom :many
SELECT * FROM events
WHERE created_at >= $1 LIMIT $2;

-- name: CreateEvent :one
INSERT INTO events
    (aggregate_id, aggregate_kind, kind, version, created_at, data)
VALUES
    ($1, $2, $3, $4, current_timestamp, $5)
RETURNING *;

-- name: GetEventsForUser :many
SELECT * FROM events
WHERE aggregate_id = $1 AND aggregate_kind = 'user' AND created_at >= $2 LIMIT $3;
