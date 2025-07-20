package handlers

import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "Chirpy/internal/auth"
    "github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
    apiKey, err := auth.GetAPIKey(r.Header)
    if err != nil {
        log.Printf("Error getting API key: %s", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    if apiKey != cfg.PolkaKey {
        log.Printf("Invalid API key provided: %s", apiKey)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    type webhookData struct {
        UserID string `json:"user_id"`
    }
    
    type webhookRequest struct {
        Event string      `json:"event"`
        Data  webhookData `json:"data"`
    }

    var req webhookRequest
    decoder := json.NewDecoder(r.Body)
    
    err = decoder.Decode(&req)
    if err != nil {
        log.Printf("Error decoding webhook request: %s", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if req.Event != "user.upgraded" {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    userID, err := uuid.Parse(req.Data.UserID)
    if err != nil {
        log.Printf("Invalid user ID in webhook: %s", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    _, err = cfg.DbQueries.UpdateUserToChirpyRed(context.Background(), userID)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        log.Printf("Error upgrading user to Chirpy Red: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}