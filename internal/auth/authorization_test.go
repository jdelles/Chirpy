package auth

import (
    "encoding/hex"
    "net/http"
    "testing"
)

func TestGetBearerToken(t *testing.T) {
    tests := []struct {
        name           string
        authHeader     string
        expectedToken  string
        expectedError  bool
    }{
        {
            name:          "valid bearer token",
            authHeader:    "Bearer abc123token",
            expectedToken: "abc123token",
            expectedError: false,
        },
        {
            name:          "valid bearer token with extra spaces",
            authHeader:    "Bearer   abc123token   ",
            expectedToken: "abc123token",
            expectedError: false,
        },
        {
            name:          "missing authorization header",
            authHeader:    "",
            expectedToken: "",
            expectedError: true,
        },
        {
            name:          "invalid format - no bearer",
            authHeader:    "abc123token",
            expectedToken: "",
            expectedError: true,
        },
        {
            name:          "invalid format - wrong prefix",
            authHeader:    "Basic abc123token",
            expectedToken: "",
            expectedError: true,
        },
        {
            name:          "bearer with no token",
            authHeader:    "Bearer",
            expectedToken: "",
            expectedError: true,
        },
        {
            name:          "bearer with empty token",
            authHeader:    "Bearer   ",
            expectedToken: "",
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            headers := make(http.Header)
            if tt.authHeader != "" {
                headers.Set("Authorization", tt.authHeader)
            }

            token, err := GetBearerToken(headers)

            if tt.expectedError && err == nil {
                t.Errorf("Expected error but got none")
            }
            if !tt.expectedError && err != nil {
                t.Errorf("Expected no error but got: %v", err)
            }
            if token != tt.expectedToken {
                t.Errorf("Expected token '%s', got '%s'", tt.expectedToken, token)
            }
        })
    }
}


func TestMakeRefreshToken(t *testing.T) {
    token, err := MakeRefreshToken()
    if err != nil {
        t.Fatalf("MakeRefreshToken failed: %v", err)
    }

    if token == "" {
        t.Fatal("MakeRefreshToken returned empty token")
    }

    if len(token) != 64 {
        t.Errorf("Expected token length 64, got %d", len(token))
    }

    _, err = hex.DecodeString(token)
    if err != nil {
        t.Errorf("Token is not valid hex: %v", err)
    }

    token2, err := MakeRefreshToken()
    if err != nil {
        t.Fatalf("Second MakeRefreshToken failed: %v", err)
    }

    if token == token2 {
        t.Error("Two consecutive tokens should be different")
    }
}

func TestGetAPIKey(t *testing.T) {
    tests := []struct {
        name           string
        authHeader     string
        expectedKey    string
        expectedError  bool
    }{
        {
            name:          "valid API key",
            authHeader:    "ApiKey abc123key",
            expectedKey:   "abc123key",
            expectedError: false,
        },
        {
            name:          "valid API key with extra spaces",
            authHeader:    "ApiKey   abc123key   ",
            expectedKey:   "abc123key",
            expectedError: false,
        },
        {
            name:          "case insensitive ApiKey",
            authHeader:    "apikey abc123key",
            expectedKey:   "abc123key",
            expectedError: false,
        },
        {
            name:          "missing authorization header",
            authHeader:    "",
            expectedKey:   "",
            expectedError: true,
        },
        {
            name:          "invalid format - no ApiKey prefix",
            authHeader:    "abc123key",
            expectedKey:   "",
            expectedError: true,
        },
        {
            name:          "invalid format - wrong prefix",
            authHeader:    "Bearer abc123key",
            expectedKey:   "",
            expectedError: true,
        },
        {
            name:          "ApiKey with no key",
            authHeader:    "ApiKey",
            expectedKey:   "",
            expectedError: true,
        },
        {
            name:          "ApiKey with empty key",
            authHeader:    "ApiKey   ",
            expectedKey:   "",
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            headers := make(http.Header)
            if tt.authHeader != "" {
                headers.Set("Authorization", tt.authHeader)
            }

            key, err := GetAPIKey(headers)

            if tt.expectedError && err == nil {
                t.Errorf("Expected error but got none")
            }
            if !tt.expectedError && err != nil {
                t.Errorf("Expected no error but got: %v", err)
            }
            if key != tt.expectedKey {
                t.Errorf("Expected key '%s', got '%s'", tt.expectedKey, key)
            }
        })
    }
}
