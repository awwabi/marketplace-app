package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type JWTClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	expTime := 60

	if os.Getenv("ENV") == "production" {
		expTime = 2
	}

	expirationTime := time.Now().Add(time.Duration(expTime) * time.Minute)

	claims := JWTClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New("invalid token or expired")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token is expired")
	}

	return &claims.ID, nil
}

func ExtractToken(tokenString string) string {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		return tokenString[7:]
	}

	return tokenString
}
