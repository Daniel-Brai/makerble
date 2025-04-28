package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/repository"
	"github.com/yhwbach/makerble/internal/schemas"
)


type MockPatientRepo struct {
	patients map[uuid.UUID]*models.Patient
	mu       sync.RWMutex
}

type MockUserRepo struct {
	users map[uuid.UUID]*models.User
	mu    sync.RWMutex
}

type MockTokenRepo struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

func NewMockRepoStorage() repository.RepoStorage {
	return repository.RepoStorage{
		Patients: &MockPatientRepo{patients: make(map[uuid.UUID]*models.Patient)},
		Users:    &MockUserRepo{users: make(map[uuid.UUID]*models.User)},
		Tokens:   &MockTokenRepo{tokens: make(map[string]time.Time)},
	}
}

func (m *MockPatientRepo) Create(ctx context.Context, userID uuid.UUID, patient *schemas.PatientCreate, createdAt time.Time) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New()
	dob, _ := time.Parse(patient.DateOfBirth, "2006-01-02")
	m.patients[id] = &models.Patient{
		ID:             id,
		FullName:       patient.FullName,
		Gender:         patient.Gender,
		Address:        patient.Address,
		DateOfBirth:    dob,
		Phone:          patient.Phone,
		Email:          patient.Email,
		MedicalHistory: patient.MedicalHistory,
		RegisteredBy:   userID,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}
	return id.String(), nil
}

func (m *MockPatientRepo) FindAll(ctx context.Context) ([]schemas.Patients, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var patients []schemas.Patients
	for _, p := range m.patients {
		patients = append(patients, schemas.Patients{
			Patient: p,
			RegisteredByUser: struct {
				ID       string `json:"id"`
				FullName string `json:"full_name"`
			}{
				ID:       p.RegisteredBy.String(),
				FullName: "Test User", // Mock data
			},
		})
	}
	return patients, nil
}

func (m *MockPatientRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Patient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if patient, exists := m.patients[id]; exists {
		return patient, nil
	}
	return nil, nil
}

func (m *MockPatientRepo) FindByEmail(ctx context.Context, email string) (*models.Patient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.patients {
		if p.Email == email {
			return p, nil
		}
	}
	return nil, nil
}

func (m *MockPatientRepo) UpdateByID(ctx context.Context, id uuid.UUID, update *schemas.PatientUpdate) (*models.Patient, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if patient, exists := m.patients[id]; exists {
		if update.FullName != nil {
			patient.FullName = *update.FullName
		}
		if update.Gender != nil {
			patient.Gender = *update.Gender
		}
		if update.Address != nil {
			patient.Address = *update.Address
		}
		if update.Phone != nil {
			patient.Phone = *update.Phone
		}
		if update.Email != nil {
			patient.Email = *update.Email
		}
		if update.MedicalHistory != nil {
			patient.MedicalHistory = *update.MedicalHistory
		}
		patient.UpdatedAt = time.Now()
		return patient, nil
	}
	return nil, nil
}

func (m *MockPatientRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.patients, id)
	return nil
}

// MockUserRepo implementations
func (m *MockUserRepo) Create(ctx context.Context, user *schemas.UserRegister, hashedPassword string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New()
	m.users[id] = &models.User{
		ID:        id,
		Username:  user.Username,
		Password:  hashedPassword,
		Email:     user.Email,
		FullName:  user.FullName,
		UserType:  user.UserType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return id.String(), nil
}

func (m *MockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepo) UpdateByID(ctx context.Context, id uuid.UUID, update *schemas.UserUpdate) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user, exists := m.users[id]; exists {
		if update.Username != nil {
			user.Username = *update.Username
		}
		if update.Email != nil {
			user.Email = *update.Email
		}
		if update.FullName != nil {
			user.FullName = *update.FullName
		}
		if update.UserType != nil {
			user.UserType = *update.UserType
		}
		user.UpdatedAt = time.Now()
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockUserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Username == username {
			return true, nil
		}
	}
	return false, nil
}

// MockTokenRepo implementations
func (m *MockTokenRepo) InvalidateToken(ctx context.Context, token string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tokens[token] = expiresAt
	return nil
}

func (m *MockTokenRepo) IsTokenInvalid(ctx context.Context, token string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if expiry, exists := m.tokens[token]; exists {
		return time.Now().Before(expiry), nil
	}
	return false, nil
}

func (m *MockTokenRepo) CleanupExpiredTokens(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, expiry := range m.tokens {
		if now.After(expiry) {
			delete(m.tokens, token)
		}
	}
	return nil
}
