package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/schemas"
	"github.com/yhwbach/makerble/internal/testutils"
)

func TestLoginHandler(t *testing.T) {
	ts := testutils.NewTestServer(t)
	defer ts.Close()

	// First register a test user
	registerUser := schemas.UserRegister{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
		FullName: "Test User",
		UserType: models.Doctor,
	}

	resp := testutils.MakeRequest(t, ts, http.MethodPost, "/api/v1/register", registerUser, "")
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	tests := []struct {
		name       string
		payload    schemas.UserLogin
		wantStatus int
	}{
		{
			name: "valid credentials",
			payload: schemas.UserLogin{
				Username: "testuser",
				Password: "testpass",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid password",
			payload: schemas.UserLogin{
				Username: "testuser",
				Password: "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "user not found",
			payload: schemas.UserLogin{
				Username: "nonexistent",
				Password: "testpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := testutils.MakeRequest(t, ts, http.MethodPost, "/v1/login", tt.payload, "")
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				var tokenResp schemas.TokenResponse
				err := json.NewDecoder(resp.Body).Decode(&tokenResp)
				require.NoError(t, err)
				assert.NotEmpty(t, tokenResp.AccessToken)
				assert.Equal(t, "Bearer", tokenResp.TokenType)
				assert.Equal(t, string(models.Doctor), tokenResp.UserType)
			}
		})
	}
}
