package server

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/yhwbach/makerble/internal/models"
	"github.com/yhwbach/makerble/internal/utils"
)

func (a *Application) authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenStr, err := utils.ExtractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		// Check if token is invalidated
		isInvalid, err := a.Repo.Tokens.IsTokenInvalid(r.Context(), tokenStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error checking token validity")
			return
		}

		if isInvalid {
			respondWithError(w, http.StatusUnauthorized, "Token has been invalidated")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Application) receptionistOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, err := utils.GetUserTypeFromContext(r.Context())
		if err != nil || userType != string(models.Receptionist) {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Application) doctorOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, err := utils.GetUserTypeFromContext(r.Context())
		if err != nil || userType != string(models.Doctor) {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	})
}
