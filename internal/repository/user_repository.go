package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/schemas"
)

type UserRepoStorage struct {
	db *sql.DB
}

func (r *UserRepoStorage) Create(ctx context.Context, user *schemas.UserRegister, hashedPassword string) (string, error) {
	userModel := &models.User{
		ID:        uuid.New(),
		Username:  user.Username,
		Password:  hashedPassword,
		Email:     user.Email,
		FullName:  user.FullName,
		UserType:  user.UserType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	query := `
		INSERT INTO users (id, username, password, email, full_name, user_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(
		ctx, query,
		userModel.ID,
		userModel.Username, userModel.Password, userModel.Email,
		userModel.FullName, userModel.UserType,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id.String(), nil
}

// GetByID retrieves a user by ID
func (r *UserRepoStorage) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, password, email, full_name, user_type, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName,
		&user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepoStorage) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password, email, full_name, user_type, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName,
		&user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepoStorage) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, password, email, full_name, user_type, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.FullName,
		&user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// UpdateByID updates a user by ID
func (r *UserRepoStorage) UpdateByID(ctx context.Context, id uuid.UUID, user *schemas.UserUpdate) (*models.User, error) {
	query := `
		UPDATE users
		SET username = COALESCE($1, username),
			email = COALESCE($2, email),
			full_name = COALESCE($3, full_name),
			user_type = COALESCE($4, user_type),
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, username, email, full_name, user_type, created_at, updated_at
	`
	var updatedUser models.User
	err := r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.FullName, user.UserType, id,
	).Scan(
		&updatedUser.ID, &updatedUser.Username, &updatedUser.Email,
		&updatedUser.FullName, &updatedUser.UserType, &updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &updatedUser, nil
}


// UsernameExists checks if a username already exists
func (r *UserRepoStorage) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	
	return exists, nil
}

// EmailExists checks if an email already exists
func (r *UserRepoStorage) EmailExists(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	
	return exists, nil
}