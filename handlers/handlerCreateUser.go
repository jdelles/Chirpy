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

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

    email := params.Email
    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, err.Error())
    }

    createUserParams := database.CreateUserParams{
        Email: email,
        HashedPassword: hashedPassword,
    }

    log.Printf("Creating user with email: %s", email)

    ctx := context.Background()
    newUser, err := cfg.DbQueries.CreateUser(ctx, createUserParams)
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, "Failed to create user")
        return
    }

    type userResponse struct {
        ID          uuid.UUID `json:"id"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        Email       string    `json:"email"`
        IsChirpyRed bool      `json:"is_chirpy_red"`
    }

    response := userResponse{
        ID:          newUser.ID,
        CreatedAt:   newUser.CreatedAt,
        UpdatedAt:   newUser.UpdatedAt,
        Email:       newUser.Email,
        IsChirpyRed: newUser.IsChirpyRed,
    }

    RespondWithJSON(w, http.StatusCreated, response)
}