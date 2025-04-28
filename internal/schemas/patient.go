package schemas

import (
	"github.com/yhwbach/makerble/internal/models"
)

// PatientCreate represents a request to create a new patient
type PatientCreate struct {
	FullName       string        `json:"full_name"`
	DateOfBirth    string        `json:"date_of_birth"` // Format: YYYY-MM-DD
	Gender         models.Gender `json:"gender"`
	Address        string        `json:"address"`
	Phone          string        `json:"phone"`
	Email          string        `json:"email"`
	MedicalHistory string        `json:"medical_history"`
}

// PatientUpdate represents a request to update an existing patient
type PatientUpdate struct {
	FullName       *string        `json:"full_name,omitempty"`
	DateOfBirth    *string        `json:"date_of_birth,omitempty"`
	Gender         *models.Gender `json:"gender,omitempty"`
	Address        *string        `json:"address,omitempty"`
	Phone          *string        `json:"phone,omitempty"`
	Email          *string        `json:"email,omitempty"`
	MedicalHistory *string        `json:"medical_history,omitempty"`
}

type PaginationQuery struct {
	Page     int `json:"page" form:"page,default=1"`
	PageSize int `json:"page_size" form:"page_size,default=10"`
}

type PatientListResponse struct {
	Patients []Patients `json:"patients"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Patients struct {
	*models.Patient
	RegisteredByUser struct {
		ID       string `json:"id"`
		FullName string `json:"full_name"`
	} `json:"registered_by_user"`
}

type PatientCreateResponse struct {
	Message string `json:"message"`
	PatientID string `json:"patient_id"`
}