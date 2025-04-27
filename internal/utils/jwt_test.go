package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTManager(t *testing.T) {
	secret := "test_secret"
	expiry := time.Hour
	manager := NewJWTManager(secret, expiry)

	t.Run("GenerateToken", func(t *testing.T) {
		userID := uuid.New()
		userType := "doctor"

		token, err := manager.GenerateToken(userID, userType)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
