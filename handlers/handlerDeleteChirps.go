package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"

	"Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerDeleteChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JwtSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
	}

	chirp, err := cfg.DbQueries.ReadChirpsByID(context.Background(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		log.Printf("Error getting chirp: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp")
		return
	}

	if chirp.UserID.UUID != userID || !chirp.UserID.Valid {
		RespondWithError(w, http.StatusForbidden, "You can only delete your own chirps")
		return
	}

	err = cfg.DbQueries.DeleteChirps(context.Background(), chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}