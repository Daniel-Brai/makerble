package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"


	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yhwbach/makerble/internal/config"
	"github.com/yhwbach/makerble/internal/repository"
	"github.com/yhwbach/makerble/internal/server"
	"github.com/yhwbach/makerble/internal/utils"
)

type TestServer struct {
	App        *server.Application
	TestServer *httptest.Server
	DB         *TestDatabase
}

func NewTestServer(t *testing.T) *TestServer {
	testDB := SetupTestDB(t)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test_secret",
			Expiry: time.Hour,
		},
	}

	repo := repository.NewRepoStorage(testDB.DB)
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiry)
	app := server.NewApplication(cfg, repo, jwtManager)

	testServer := httptest.NewServer(app.Mount())

	return &TestServer{
		App:        app,
		TestServer: testServer,
		DB:         testDB,
	}
}

func (ts *TestServer) Close() {
	ts.TestServer.Close()
	ts.DB.Close()
}

func GenerateTestToken(t *testing.T, userID uuid.UUID, userType string, jwtManager *utils.JWTManager) string {
	token, err := jwtManager.GenerateToken(userID, userType)
	require.NoError(t, err)
	return token
}

func MakeRequest(t *testing.T, ts *TestServer, method, path string, body interface{}, token string) *http.Response {
	var reqBody bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&reqBody).Encode(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, ts.TestServer.URL+path, &reqBody)
	require.NoError(t, err)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
