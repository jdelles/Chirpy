-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    Now(),
    Now(),
    $1, 
    $2
)
RETURNING *;