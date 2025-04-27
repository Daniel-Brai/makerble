package models

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

// Patient represents a patient in the medical system
type Patient struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"full_name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Gender         Gender    `json:"gender"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	MedicalHistory string    `json:"medical_history"`
	RegisteredBy   uuid.UUID `json:"registered_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
