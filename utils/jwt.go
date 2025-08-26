package utils

import (
	"fmt"
	"time"
	"github.com/Okeke-Divine/flaking-api/config"
	"github.com/golang-jwt/jwt/v5"
)

// Custom claims structure
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint, email string) (string, error) {
	// DEBUG: Check if JWT secret is loaded
	fmt.Printf("JWT Secret: '%s'\n", config.AppConfig.JWTSecret)
	if config.AppConfig.JWTSecret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}
	
	// Set token expiration time (7 days)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates JWT token and returns claims
func ValidateToken(tokenString string) (*Claims, error) {
	// DEBUG: Check if JWT secret is loaded
	if config.AppConfig.JWTSecret == "" {
		return nil, fmt.Errorf("JWT secret is not configured")
	}
	
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}