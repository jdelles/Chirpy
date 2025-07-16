package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email string `json:"email"`
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

    log.Printf("Creating user with email: %s", email)

    ctx := context.Background()
    user, err := cfg.DbQueries.CreateUser(ctx, email)
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, "Failed to create user")
        return
    }

    RespondWithJSON(w, http.StatusCreated, user)
}