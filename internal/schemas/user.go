package schemas

import (
	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/models"
)

// UserUpdate represents a request to update user information
type UserUpdate struct {
	Username *string          `json:"username,omitempty"`
	Email    *string          `json:"email,omitempty"`
	FullName *string          `json:"full_name,omitempty"`
	UserType *models.UserType `json:"user_type,omitempty"`
}

// UserLogin represents login request body
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserRegister represents registration request body
type UserRegister struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Email    string          `json:"email"`
	FullName string          `json:"full_name"`
	UserType models.UserType `json:"user_type"`
}

// TokenResponse represents the response after successful authentication
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UserType    string `json:"user_type"`
}

// UserResponse represents the response for user information
type UserRegisterResponse struct {
	Message string    `json:"message"`
	UserID  uuid.UUID `json:"user_id"`
}
