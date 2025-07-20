package handlers

import (
    "context"
    "log"
    "net/http"

    "github.com/google/uuid"
    "Chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerReadChirps(w http.ResponseWriter, r *http.Request) {
    authorIDParam := r.URL.Query().Get("author_id")
    sortParam := r.URL.Query().Get("sort")
    sortDesc := (sortParam == "desc")
    
    var chirps []database.Chirp
    var err error
    
    if authorIDParam == "" {
        if sortDesc {
            chirps, err = cfg.DbQueries.ReadChirpsDesc(context.Background())
        } else {
            chirps, err = cfg.DbQueries.ReadChirps(context.Background())
        }
    } else {
        authorID, parseErr := uuid.Parse(authorIDParam)
        if parseErr != nil {
            log.Printf("Invalid author_id format: %s", parseErr)
            RespondWithError(w, http.StatusBadRequest, "Invalid author_id format")
            return
        }
        
        authorUUID := uuid.NullUUID{
            UUID:  authorID,
            Valid: true,
        }
        
        if sortDesc {
            chirps, err = cfg.DbQueries.ReadChirpsByAuthorDesc(context.Background(), authorUUID)
        } else {
            chirps, err = cfg.DbQueries.ReadChirpsByAuthor(context.Background(), authorUUID)
        }
    }
    
    if err != nil {
        log.Printf("Error getting chirps: %s", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
        return
    }
    
    RespondWithJSON(w, http.StatusOK, chirps)
}