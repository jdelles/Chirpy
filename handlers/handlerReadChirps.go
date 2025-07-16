package handlers

import (
	"context"
	"net/http"
)

func (cfg *ApiConfig) HandlerReadChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DbQueries.ReadChirps(context.Background())
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to get chirps")
		return
	}
	RespondWithJSON(w, http.StatusOK, chirps)
}