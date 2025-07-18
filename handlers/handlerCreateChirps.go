package handlers

import (
    "context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
    
	"github.com/google/uuid"
    
	"Chirpy/internal/auth"
    "Chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerCreateChirps(w http.ResponseWriter, r *http.Request) {
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
        Body string `json:"body"`
    }
    params := parameters{}

    type returnVals struct {
        Cleaned_body string `json:"cleaned_body"`
    }
    
    decoder := json.NewDecoder(r.Body)
    
    err = decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters %s", err)
        RespondWithError(w, http.StatusBadRequest, "Something went wrong")
        return
    }
    
    if len(params.Body) > 140 {
        RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    cleanedBody := cleanProfaneWords(params.Body)

	createChirpParams := database.CreateChirpParams{
		Body: cleanedBody,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	}

	chirp, err := cfg.DbQueries.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to create chirp")
	}

    RespondWithJSON(w, http.StatusCreated, chirp)
}

func cleanProfaneWords(body string) string {
    replacement := "****"
    badWords := []string{"kerfuffle", "sharbert", "fornax"}
    splitString := strings.Split(body, " ")
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
    
    return strings.Join(cleanedWords, " ")
}
