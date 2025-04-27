package models

import (
	"time"

	"github.com/google/uuid"
)

type InvalidToken struct {
	ID            uuid.UUID `json:"id"`
	Token         string    `json:"token"`
	InvalidatedAt time.Time `json:"invalidated_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}
