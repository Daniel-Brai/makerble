package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/schemas"
	"github.com/yhwbach/makerble/internal/testutils"
)

func TestCreatePatientHandler(t *testing.T) {
	ts := testutils.NewTestServer(t)
	defer ts.Close()

	// Create a test user first
	userID := uuid.New()
	token := testutils.GenerateTestToken(t, userID, string(models.Receptionist), ts.App.JWTManager)

	tests := []struct {
		name         string
		body         schemas.PatientCreate
		token        string
		expectedCode int
	}{
		{
			name: "valid patient creation",
			body: schemas.PatientCreate{
				FullName:    "John Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0).Format("2006-01-02"),
				Gender:      models.Male,
				Email:       "john@example.com",
			},
			token:        token,
			expectedCode: http.StatusCreated,
		},
		{
			name: "unauthorized access",
			body: schemas.PatientCreate{
				FullName:    "John Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0).Format("2006-01-02"),
				Gender:      models.Male,
				Email:       "john@example.com",
			},
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := testutils.MakeRequest(t, ts, "POST", "/v1/patients", tt.body, tt.token)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedCode == http.StatusCreated {
				var response map[string]string
				err := json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.NotEmpty(t, response["id"])
			}
		})
	}
}

