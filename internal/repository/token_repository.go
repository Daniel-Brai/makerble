package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TokenRepoStorage struct {
	db *sql.DB
}

func (r *TokenRepoStorage) InvalidateToken(ctx context.Context, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO invalid_tokens (token, expires_at)
		VALUES ($1, $2)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, token, expiresAt).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %w", err)
	}

	return nil
}

func (r *TokenRepoStorage) IsTokenInvalid(ctx context.Context, token string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM invalid_tokens 
			WHERE token = $1 
			AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check token validity: %w", err)
	}

	return exists, nil
}

func (r *TokenRepoStorage) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM invalid_tokens WHERE expires_at <= NOW()`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}
