-- name: ReadUserByEmail :one
SELECT *
FROM users
WHERE email = $1;