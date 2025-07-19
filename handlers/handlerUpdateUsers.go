package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"Chirpy/internal/auth"
	"Chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
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

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Email == "" || params.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	updateUserParams := database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	updatedUser, err := cfg.DbQueries.UpdateUser(context.Background(), updateUserParams)
	if err != nil {
		log.Printf("Error updating user: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

    type userResponse struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email     string    `json:"email"`
    }

    response := userResponse{
        ID:        updatedUser.ID,
        CreatedAt: updatedUser.CreatedAt,
        UpdatedAt: updatedUser.UpdatedAt,
        Email:     updatedUser.Email,
    }

    RespondWithJSON(w, http.StatusOK, response)
}