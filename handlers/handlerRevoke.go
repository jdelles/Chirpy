package handlers

import (
    "context"
    "log"
    "net/http"

	"Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
    refreshToken, err := auth.GetBearerToken(r.Header)
    if err != nil {
        log.Printf("Error getting bearer token: %v", err)
        RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }
    
    log.Printf("Attempting to revoke token: %s", refreshToken)
    
    ctx := context.Background()
    err = cfg.DbQueries.RevokeRefreshToken(ctx, refreshToken)
    if err != nil {
        log.Printf("Error revoking refresh token: %v", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
        return
    }
    
    log.Printf("Successfully revoked token")
    w.WriteHeader(http.StatusNoContent)
}