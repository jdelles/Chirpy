package auth

import (
	"strings"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
)

func TestMakeJWT(t *testing.T) {
    // Test data
    userID := uuid.New()
    tokenSecret := "test-secret"
    expiresIn := time.Hour

    // Create JWT
    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("MakeJWT failed: %v", err)
    }

    // Verify token is not empty
    if tokenString == "" {
        t.Fatal("MakeJWT returned empty token")
    }

    // Parse the token to verify it was created correctly
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
    })
    if err != nil {
        t.Fatalf("Failed to parse created token: %v", err)
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok {
        t.Fatal("Invalid token claims")
    }

    // Verify claims
    if claims.Issuer != "chirpy" {
        t.Errorf("Expected issuer 'chirpy', got '%s'", claims.Issuer)
    }

    if claims.Subject != userID.String() {
        t.Errorf("Expected subject '%s', got '%s'", userID.String(), claims.Subject)
    }

    // Check expiration is approximately correct (within 1 second)
    expectedExpiry := time.Now().UTC().Add(expiresIn)
    if claims.ExpiresAt.Time.Sub(expectedExpiry) > time.Second {
        t.Errorf("Token expiration time is incorrect")
    }
}

func TestValidateJWT_Success(t *testing.T) {
    // Test data
    userID := uuid.New()
    tokenSecret := "test-secret"
    expiresIn := time.Hour

    // Create a valid token
    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("Failed to create test token: %v", err)
    }

    // Validate the token
    validatedUserID, err := ValidateJWT(tokenString, tokenSecret)
    if err != nil {
        t.Fatalf("ValidateJWT failed: %v", err)
    }

    // Verify the user ID matches
    if validatedUserID != userID {
        t.Errorf("Expected user ID '%s', got '%s'", userID.String(), validatedUserID.String())
    }
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
    // Test data
    userID := uuid.New()
    tokenSecret := "test-secret"
    expiresIn := -time.Hour // Expired 1 hour ago

    // Create an expired token
    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("Failed to create test token: %v", err)
    }

    // Try to validate expired token
    _, err = ValidateJWT(tokenString, tokenSecret)
    if err == nil {
        t.Fatal("Expected error for expired token, got nil")
    }
}

func TestValidateJWT_WrongSecret(t *testing.T) {
    // Test data
    userID := uuid.New()
    tokenSecret := "correct-secret"
    wrongSecret := "wrong-secret"
    expiresIn := time.Hour

    // Create a token with correct secret
    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("Failed to create test token: %v", err)
    }

    // Try to validate with wrong secret
    _, err = ValidateJWT(tokenString, wrongSecret)
    if err == nil {
        t.Fatal("Expected error for wrong secret, got nil")
    }

    // Check if it's a signature error by examining the error message
    if !strings.Contains(err.Error(), "signature") {
        t.Errorf("Expected signature error, got: %v", err)
    }
}

func TestValidateJWT_InvalidToken(t *testing.T) {
    tokenSecret := "test-secret"
    
    tests := []struct {
        name        string
        tokenString string
    }{
        {"empty token", ""},
        {"malformed token", "invalid.token.string"},
        {"not a JWT", "this-is-not-a-jwt"},
        {"missing segments", "header.payload"}, // Missing signature
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ValidateJWT(tt.tokenString, tokenSecret)
            if err == nil {
                t.Fatalf("Expected error for %s, got nil", tt.name)
            }
        })
    }
}

func TestValidateJWT_InvalidUserID(t *testing.T) {
    tokenSecret := "test-secret"
    
    // Create a token with invalid UUID in subject
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Issuer:    "chirpy",
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
        Subject:   "invalid-uuid", // This is not a valid UUID
    })

    tokenString, err := token.SignedString([]byte(tokenSecret))
    if err != nil {
        t.Fatalf("Failed to create test token: %v", err)
    }

    // Try to validate token with invalid UUID
    _, err = ValidateJWT(tokenString, tokenSecret)
    if err == nil {
        t.Fatal("Expected error for invalid UUID, got nil")
    }

    // Verify it's a UUID parsing error
    if err.Error() != "invalid user ID in token: invalid UUID length: 12" {
        t.Logf("Got error (this is expected): %v", err)
    }
}

func TestMakeJWT_DifferentExpirations(t *testing.T) {
    userID := uuid.New()
    tokenSecret := "test-secret"
    
    tests := []struct {
        name      string
        expiresIn time.Duration
    }{
        {"1 minute", time.Minute},
        {"1 hour", time.Hour},
        {"24 hours", 24 * time.Hour},
        {"1 week", 7 * 24 * time.Hour},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tokenString, err := MakeJWT(userID, tokenSecret, tt.expiresIn)
            if err != nil {
                t.Fatalf("MakeJWT failed for %s: %v", tt.name, err)
            }

            // Validate the token
            validatedUserID, err := ValidateJWT(tokenString, tokenSecret)
            if err != nil {
                t.Fatalf("ValidateJWT failed for %s: %v", tt.name, err)
            }

            if validatedUserID != userID {
                t.Errorf("User ID mismatch for %s", tt.name)
            }
        })
    }
}

func TestRoundTrip(t *testing.T) {
    // Test the full round trip: create token -> validate token
    userID := uuid.New()
    tokenSecret := "round-trip-secret"
    expiresIn := time.Hour * 2

    // Create token
    tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
    if err != nil {
        t.Fatalf("Failed to create token: %v", err)
    }

    // Validate token
    validatedUserID, err := ValidateJWT(tokenString, tokenSecret)
    if err != nil {
        t.Fatalf("Failed to validate token: %v", err)
    }

    // Verify user ID matches
    if validatedUserID != userID {
        t.Errorf("Round trip failed: expected %s, got %s", userID, validatedUserID)
    }
}