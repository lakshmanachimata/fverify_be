package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_secret_key") // Replace with a secure secret key

type AuthTokenClaims struct {
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	Status       string `json:"status"`
	MobileNumber string `json:"mobileNumber"`
	jwt.RegisteredClaims
}

// GenerateAuthToken generates a JWT token for the user
func GenerateAuthToken(userId, username, role, status, mobileNumber string) (string, error) {
	claims := AuthTokenClaims{
		UserId:       userId,
		Username:     username,
		Role:         role,
		Status:       status,
		MobileNumber: mobileNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseAuthToken parses and validates a JWT token
func ParseAuthToken(tokenString string) (*AuthTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
