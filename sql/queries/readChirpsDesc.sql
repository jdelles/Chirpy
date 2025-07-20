-- name: ReadChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;