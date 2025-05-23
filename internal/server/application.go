package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/yhwbach/makerble/docs"
	"github.com/yhwbach/makerble/internal/config"
	"github.com/yhwbach/makerble/internal/repository"
	"github.com/yhwbach/makerble/internal/utils"
)

type Application struct {
	Config     *config.Config
	Repo       repository.RepoStorage
	JWTManager *utils.JWTManager
}

func NewApplication(cfg *config.Config, repo repository.RepoStorage, jwtManager *utils.JWTManager) *Application {
	return &Application{
		Config:     cfg,
		Repo:       repo,
		JWTManager: jwtManager,
	}
}

func GetProjectRoot() (string, error) {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return "", fmt.Errorf("failed to get current file path")
    }
    
    dir := filepath.Dir(filename)
    
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir, nil
        }
        
        parent := filepath.Dir(dir)
        if parent == dir {
            return "", fmt.Errorf("go.mod not found")
        }
        dir = parent
    }
}

func (a *Application) Mount() http.Handler {
	r := chi.NewRouter()

	prefix := "/api/v1"

	// Basic middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route(prefix, func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(middleware.Compress(5, "application/json"))

		// Health check
		r.Get("/health", a.healthCheck)

		// Auth routes (no authentication required)
		r.Post("/register", a.registerHandler)
		r.Post("/login", a.loginHandler)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(a.JWTManager.Auth))
			r.Use(a.authenticator)

			r.Post("/logout", a.logoutHandler)

			// Patient routes
			r.Route("/patients", func(r chi.Router) {
				r.Get("/", a.listPatientsHandler)
				r.Get("/{id}", a.getPatientHandler)

				// Receptionist only routes
				r.Group(func(r chi.Router) {
					r.Use(a.receptionistOnly)
					r.Put("/{id}", a.updatePatientLimitedHandler)
					r.Delete("/{id}", a.deletePatientHandler)
				})

				// Doctor only routes
				r.Group(func(r chi.Router) {
					r.Use(a.doctorOnly)
					r.Post("/", a.createPatientHandler)
					r.Patch("/{id}", a.updatePatientHandler)
				})
			})
		})

	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/docs/swagger.json"),
	))

	return r
}

// Add health check handler
func (a *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.Config.Server.Port),
		Handler:      a.Mount(),
		WriteTimeout: a.Config.Server.WriteTimeout,
		ReadTimeout:  a.Config.Server.ReadTimeout,
		IdleTimeout:  a.Config.Server.IdleTimeout,
	}

	return srv.ListenAndServe()
}
