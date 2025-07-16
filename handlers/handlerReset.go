package handlers

import (
	"context"
	"net/http"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
        RespondWithError(w, http.StatusForbidden, "Reset is forbidden")
    }
    err := cfg.DbQueries.DeleteUsers(context.Background())
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, "Reset request unsuccessful")
    }
    cfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}