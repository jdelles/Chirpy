-- name: ReadChirps :many
SELECT * 
FROM chirps
ORDER BY created_at ASC;