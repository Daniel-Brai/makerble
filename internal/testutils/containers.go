package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	Container *postgres.PostgresContainer
	DB        *sql.DB
}

func SetupTestDB(t *testing.T) *TestDatabase {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		t.Fatal(err)
	}

	return &TestDatabase{
		Container: pgContainer,
		DB:        db,
	}
}

func (td *TestDatabase) Close() error {
	if err := td.DB.Close(); err != nil {
		return fmt.Errorf("error closing db: %w", err)
	}
	return td.Container.Terminate(context.Background())
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		"000001_create_users_table.up.sql",
		"000002_create_patients_table.up.sql",
		"000003_create_invalid_tokens_table.up.sql",
	}

	for _, migration := range migrations {
		content, err := os.ReadFile(fmt.Sprintf("../../migrations/%s", migration))
		if err != nil {
			return fmt.Errorf("error reading migration %s: %w", migration, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("error executing migration %s: %w", migration, err)
		}
	}

	return nil
}
