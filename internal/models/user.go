package models

import (
	"time"

	"github.com/google/uuid"
)

// UserType represents user roles in the system
type UserType string

const (
	Doctor       UserType = "doctor"
	Receptionist UserType = "receptionist"
)

// User represents a user in the system (doctor or receptionist)
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose the password
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	UserType  UserType  `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
