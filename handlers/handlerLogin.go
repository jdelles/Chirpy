package handlers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "Chirpy/internal/auth"
    "Chirpy/internal/database"
    "github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    params := parameters{}
    
    decoder := json.NewDecoder(r.Body)
    
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters %s", err)
        RespondWithError(w, http.StatusBadRequest, "Something went wrong")
        return
    }

    user, err := cfg.DbQueries.ReadUserByEmail(context.Background(), params.Email)
    if err != nil {
        RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
        return
    }

    err = auth.CheckPassword(params.Password, user.HashedPassword)
    if err != nil {
        RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
        return
    }

    accessToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, time.Hour)
    if err != nil {
        log.Printf("Error creating JWT: %s", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to create token")
        return
    }

    refreshToken, err := auth.MakeRefreshToken()
    if err != nil {
        log.Printf("Error creating refresh token: %s", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to create refresh token")
        return
    }

    _, err = cfg.DbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour), // 60 days
    })
    if err != nil {
        log.Printf("Error storing refresh token: %s", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to store refresh token")
        return
    }

type loginResponse struct {
    ID           uuid.UUID `json:"id"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    Email        string    `json:"email"`
    IsChirpyRed  bool      `json:"is_chirpy_red"`
    Token        string    `json:"token"`
    RefreshToken string    `json:"refresh_token"`
}

response := loginResponse{
    ID:           user.ID,
    CreatedAt:    user.CreatedAt,
    UpdatedAt:    user.UpdatedAt,
    Email:        user.Email,
    IsChirpyRed:  user.IsChirpyRed,
    Token:        accessToken,
    RefreshToken: refreshToken,
}

    RespondWithJSON(w, http.StatusOK, response)
}