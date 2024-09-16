package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("your-secret-key") // Must match the Authorization Server's secret

// ValidateJWT parses and validates a JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// CreateJWT generates a JWT token for a given user ID
func CreateJWT(userID string) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"sub":   userID,
		"exp":   time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
		"iat":   time.Now().Unix(),
		"scope": "read",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
