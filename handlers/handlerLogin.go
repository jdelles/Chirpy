package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	 type parameters struct {
        Email string `json:"email"`
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
	}

	err = auth.CheckPassword(params.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
	}

	type userResponse struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email     string    `json:"email"`
    }

    response := userResponse{
        ID:        user.ID,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        Email:     user.Email,
    }

    RespondWithJSON(w, http.StatusOK, response)
}