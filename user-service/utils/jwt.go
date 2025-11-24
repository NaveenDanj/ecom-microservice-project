package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": "user",
		"exp":  jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
