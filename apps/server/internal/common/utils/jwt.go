package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func GenerateToken(payload TokenPayload, jwtExpiryTime string) (string, error) {

	expiresAtTime, err := time.ParseDuration(jwtExpiryTime)
	if err != nil {
		return "", err
	}
	payload.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "coderz.space",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAtTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
