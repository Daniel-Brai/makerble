package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/schemas"
)

// RepoStorage is a struct that holds references to different repository interfaces.
type RepoStorage struct {
	Patients PatientRepository
	Users    UserRepository
	Tokens   TokenRepository
}

// PatientRepoStorage is a struct that implements the PatientRepository interface.
type PatientRepository interface {
	Create(context.Context, uuid.UUID, *schemas.PatientCreate) (string, error)
	FindAll(context.Context) ([]schemas.Patients, error) // Changed return type
	FindByID(context.Context, uuid.UUID) (*models.Patient, error)
	FindByEmail(context.Context, string) (*models.Patient, error)
	UpdateByID(context.Context, uuid.UUID, *schemas.PatientUpdate) (*models.Patient, error)
	DeleteByID(context.Context, uuid.UUID) error
}

// UserRepoStorage is a struct that implements the UserRepository interface.
type UserRepository interface {
	Create(context.Context, *schemas.UserRegister, string) (string, error)
	FindByID(context.Context, uuid.UUID) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	UpdateByID(context.Context, uuid.UUID, *schemas.UserUpdate) (*models.User, error)
	EmailExists(context.Context, string) (bool, error)
	UsernameExists(context.Context, string) (bool, error)
}

type TokenRepository interface {
	InvalidateToken(context.Context, string, time.Time) error
	IsTokenInvalid(context.Context, string) (bool, error)
	CleanupExpiredTokens(context.Context) error
}

func NewRepoStorage(db *sql.DB) RepoStorage {
	return RepoStorage{
		Patients: &PatientRepoStorage{db: db},
		Users:    &UserRepoStorage{db: db},
		Tokens:   &TokenRepoStorage{db: db},
	}
}
