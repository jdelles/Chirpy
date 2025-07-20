package handlers

import (
	"sync/atomic"

	"Chirpy/internal/database"
)

type ApiConfig struct {
    DbQueries      *database.Queries
	FileserverHits atomic.Int32
    Platform       string
	JwtSecret      string
	PolkaKey       string
}