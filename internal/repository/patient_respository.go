package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/schemas"
)

type PatientRepoStorage struct {
	db *sql.DB
}

// Create inserts a new patient into the database.
func (p *PatientRepoStorage) Create(ctx context.Context, userID uuid.UUID, patient *schemas.PatientCreate, dateOfBirth time.Time) (string, error) {
	patientModel := &models.Patient{
		FullName:       patient.FullName,
		DateOfBirth:    dateOfBirth,
		Gender:         patient.Gender,
		Address:        patient.Address,
		Phone:          patient.Phone,
		Email:          patient.Email,
		MedicalHistory: patient.MedicalHistory,
		RegisteredBy:   userID,
	}

	query := `INSERT INTO patients (full_name, date_of_birth, gender, address, phone, email, medical_history, registered_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := p.db.QueryRowContext(ctx, query,
		patientModel.FullName,
		patientModel.DateOfBirth,
		patientModel.Gender,
		patientModel.Address,
		patientModel.Phone,
		patientModel.Email,
		patientModel.MedicalHistory,
		patientModel.RegisteredBy,
	).Scan(&patientModel.ID)

	if err != nil {
		return "", err
	}

	return patientModel.ID.String(), nil
}

// FindAll retrieves all patients with their registered user details from the database.
func (p *PatientRepoStorage) FindAll(ctx context.Context) ([]schemas.Patients, error) {
	query := `
		SELECT 
			p.id, p.full_name, p.date_of_birth, p.gender, p.address, 
			p.phone, p.email, p.medical_history, p.registered_by, 
			p.created_at, p.updated_at,
			u.id as user_id, u.full_name as user_full_name
		FROM patients p
		JOIN users u ON p.registered_by = u.id
		ORDER BY p.created_at DESC
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []schemas.Patients
	for rows.Next() {
		var patient models.Patient
		var user struct {
			ID       string
			FullName string
		}

		if err := rows.Scan(
			&patient.ID,
			&patient.FullName,
			&patient.DateOfBirth,
			&patient.Gender,
			&patient.Address,
			&patient.Phone,
			&patient.Email,
			&patient.MedicalHistory,
			&patient.RegisteredBy,
			&patient.CreatedAt,
			&patient.UpdatedAt,
			&user.ID,
			&user.FullName,
		); err != nil {
			return nil, err
		}

		patientWithUser := schemas.Patients{
			Patient: &patient,
			RegisteredByUser: struct {
				ID       string `json:"id"`
				FullName string `json:"full_name"`
			}{
				ID:       user.ID,
				FullName: user.FullName,
			},
		}

		patients = append(patients, patientWithUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

// FindByID retrieves a patient by ID from the database.
func (p *PatientRepoStorage) FindByID(ctx context.Context, id uuid.UUID) (*models.Patient, error) {
	query := `SELECT id, full_name, date_of_birth, gender, address, phone, email, medical_history, registered_by, created_at, updated_at
		FROM patients WHERE id = $1`

	row := p.db.QueryRowContext(ctx, query, id)

	var patient models.Patient
	if err := row.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.Address,
		&patient.Phone,
		&patient.Email,
		&patient.MedicalHistory,
		&patient.RegisteredBy,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &patient, nil
}

// FindByEmail retrieves a patient by email from the database.
func (p *PatientRepoStorage) FindByEmail(ctx context.Context, email string) (*models.Patient, error) {
	query := `SELECT id, full_name, date_of_birth, gender, address, phone, email, medical_history, registered_by, created_at, updated_at
		FROM patients WHERE email = $1`

	row := p.db.QueryRowContext(ctx, query, email)

	var patient models.Patient
	if err := row.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.Address,
		&patient.Phone,
		&patient.Email,
		&patient.MedicalHistory,
		&patient.RegisteredBy,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &patient, nil
}

// UpdateByID updates a patient's information in the database.
func (p *PatientRepoStorage) UpdateByID(ctx context.Context, id uuid.UUID, patient *schemas.PatientUpdate) (*models.Patient, error) {
	query := `
		UPDATE patients 
		SET full_name = COALESCE($1, full_name), 
		date_of_birth = COALESCE($2, date_of_birth),
		gender = COALESCE($3, gender), 
		address = COALESCE($4, address), 
		phone = COALESCE($5, phone),
		email = COALESCE($6, email), 
		medical_history = COALESCE($7, medical_history),
		updated_at = NOW()
		WHERE id = $8 
		RETURNING id, full_name, date_of_birth, gender, address, phone, email, medical_history, registered_by, created_at, updated_at
		`

	var patientModel models.Patient
	err := p.db.QueryRowContext(ctx, query,
		patient.FullName,
		patient.DateOfBirth,
		patient.Gender,
		patient.Address,
		patient.Phone,
		patient.Email,
		patient.MedicalHistory,
		id,
	).Scan(
		&patientModel.ID,
		&patientModel.FullName,
		&patientModel.DateOfBirth,
		&patientModel.Gender,
		&patientModel.Address,
		&patientModel.Phone,
		&patientModel.Email,
		&patientModel.MedicalHistory,
		&patientModel.RegisteredBy,
		&patientModel.CreatedAt,
		&patientModel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &patientModel, nil
}

// DeleteByID deletes a patient by ID from the database.
func (p *PatientRepoStorage) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM patients WHERE id = $1`

	_, err := p.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}
	return nil
}
