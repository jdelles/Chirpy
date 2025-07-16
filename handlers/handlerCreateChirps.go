package handlers

import (
	"Chirpy/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerCreateChirps(w http.ResponseWriter, r *http.Request) {
 type parameters struct {
        Body string `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
    }
    params := parameters{}

    type returnVals struct {
        Cleaned_body string `json:"cleaned_body"`
    }
    
    decoder := json.NewDecoder(r.Body)
    
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters %s", err)
        RespondWithError(w, http.StatusBadRequest, "Something went wrong")
        return
    }
    
    if len(params.Body) > 140 {
        RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    replacement := "****"
    badWords := []string{"kerfuffle", "sharbert", "fornax"}
    splitString := strings.Split(params.Body, " ")
    cleanedWords := make([]string, len(splitString))
    for i, word := range splitString {
        found := false
        for _, bad := range badWords {
            if strings.ToLower(word) == bad {
                found = true
                break
            }
        }
        if found {
            cleanedWords[i] = replacement
        } else {
            cleanedWords[i] = word
        }
    }
    result := strings.Join(cleanedWords, " ")

	createChirpParams := database.CreateChirpParams{
		Body: result,
		UserID: params.UserID,
	}

	chirp, err := cfg.DbQueries.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to create chirp")
	}

    RespondWithJSON(w, http.StatusCreated, chirp)
}