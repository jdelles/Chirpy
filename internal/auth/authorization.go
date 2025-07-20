package auth

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "net/http"
    "strings"
)

func GetBearerToken(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", fmt.Errorf("no authorization header found")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        return "", fmt.Errorf("invalid authorization header format")
    }

    token := strings.TrimSpace(parts[1])
    if token == "" {
        return "", fmt.Errorf("no token found in authorization header")
    }

    return token, nil
}

func MakeRefreshToken() (string, error) {
    randomBytes := make([]byte, 32)
    

    _, err := rand.Read(randomBytes)
    if err != nil {
        return "", fmt.Errorf("failed to generate random bytes: %w", err)
    }
    
    token := hex.EncodeToString(randomBytes)
    
    return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", fmt.Errorf("no authorization header found")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "apikey" {
        return "", fmt.Errorf("invalid authorization header format")
    }

    apiKey := strings.TrimSpace(parts[1])
    if apiKey == "" {
        return "", fmt.Errorf("no API key found in authorization header")
    }

    return apiKey, nil
}