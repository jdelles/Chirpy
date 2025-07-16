-- name: ReadChirpsByID :one
SELECT *
FROM chirps
WHERE id = $1;