package handlers

import (
    "context"
    "log"
    "net/http"
    "time"

    "Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
    refreshToken, err := auth.GetBearerToken(r.Header)
    if err != nil {
        RespondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header")
        return
    }

    user, err := cfg.DbQueries.GetUserFromRefreshToken(context.Background(), refreshToken)
    if err != nil {
        RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
        return
    }

    accessToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, time.Hour)
    if err != nil {
        log.Printf("Error creating JWT: %s", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to create token")
        return
    }

    type refreshResponse struct {
        Token string `json:"token"`
    }

    response := refreshResponse{
        Token: accessToken,
    }

    RespondWithJSON(w, http.StatusOK, response)
}