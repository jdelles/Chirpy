package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", fmt.Errorf("you must supply a password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error while setting password")
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hash string) (error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}