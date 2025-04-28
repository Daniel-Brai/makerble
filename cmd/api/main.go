package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yhwbach/makerble/internal/config"
	"github.com/yhwbach/makerble/internal/database"
	"github.com/yhwbach/makerble/internal/repository"
	"github.com/yhwbach/makerble/internal/server"
	"github.com/yhwbach/makerble/internal/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title           Makerble Medical System API
// @version         1.0
// @description     API for managing a medical clinic system with doctors, receptionists, and patients
// @termsOfService  http://swagger.io/terms/

// @host      localhost:5000
// @BasePath  /api/v1
// @securityDefinitions.apikey  BearerAuth
// @in                         header
// @name                       Authorization
// @description               Bearer token authentication

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := database.New(
		cfg.Database.DatabaseURL(),
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.MaxIdleTime,
	)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Panic(err)
	}

	root, err := server.GetProjectRoot()
	if err != nil {
		log.Panic("failed to get project root:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file:%s/migrations", root), "postgres", driver)
	if err != nil {
		log.Panic("migration setup error %w ", err.Error())
	}

	err = m.Up()
	if err != nil {

		if err != migrate.ErrNoChange {
			log.Printf("migration error %s ", err.Error())
		} else {
			log.Printf("database is already up to date")
		}
	}

	repo := repository.NewRepoStorage(db)
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiry)

	app := server.NewApplication(cfg, repo, jwtManager)

	// Start token cleanup job
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			if err := repo.Tokens.CleanupExpiredTokens(context.Background()); err != nil {
				log.Printf("Error cleaning up expired tokens: %v", err)
			}
		}
	}()

	log.Printf("server is running at %s", cfg.Server.Port)
	if err := app.Run(); err != nil {
		log.Fatal("server error:", err)
	}
}
