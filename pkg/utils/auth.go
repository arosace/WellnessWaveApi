package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your_secret_key")

func GenerateVerificationToken(email string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func DecodeJWT(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading jwt: %w", err)
	}

	if !t.Valid {
		return nil, errors.New("JWT is not valid")
	}

	tokenMap := make(map[string]interface{})
	for key, val := range claims {
		tokenMap[key] = val
	}

	return tokenMap, nil
}
