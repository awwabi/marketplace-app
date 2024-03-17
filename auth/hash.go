package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

func HashPassword(password string) (string, error) {
	saltEnv := os.Getenv("BCRYPT_SALT")
	if saltEnv == "" {
		return "", fmt.Errorf("BCRYPT_SALT environment variable is not set")
	}

	salt, err := strconv.Atoi(saltEnv)
	if err != nil {
		return "", fmt.Errorf("Error parsing BCRYPT_SALT: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	if err != nil {
		return "", fmt.Errorf("Error generating bcrypt hash: %v", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) bool {
	// Compare the hashed password with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
