package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
    type errorResponse struct {
        Error string `json:"error"`
    }
    RespondWithJSON(w, code, errorResponse{Error: msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshaling JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Write(dat)
}