package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerReadChirpsByID(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to get chirp by id")
	}
	chirpsByID, err := cfg.DbQueries.ReadChirpsByID(context.Background(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Unable to get chirp by id")
	}
	RespondWithJSON(w, http.StatusOK, chirpsByID)
}