package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	Secret string
	Expiry time.Duration
	Auth   *jwtauth.JWTAuth
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		Secret: secret,
		Expiry: expiry,
		Auth:   jwtauth.New("HS256", []byte(secret), nil),
	}
}

// GenerateToken generates a new JWT token for a user
func (m *JWTManager) GenerateToken(userID uuid.UUID, userType string) (string, error) {
	claims := map[string]interface{}{
		"user_id":   userID.String(),
		"user_type": userType,
		"exp":       time.Now().Add(m.Expiry).Unix(),
	}

	_, tokenString, err := m.Auth.Encode(claims)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// GetUserIDFromContext extracts the user ID from the JWT claims in the context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user_id in token")
	}

	return userID, nil
}

// GetUserTypeFromContext extracts the user type from the JWT claims in the context
func GetUserTypeFromContext(ctx context.Context) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", err
	}

	userType, ok := claims["user_type"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user_type in token")
	}

	return userType, nil
}

// ExtractBearerToken extracts the token from the Authorization header
func ExtractBearerToken(authHeader string) (string, error) {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}
	return "", fmt.Errorf("invalid authorization header")
}
