package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/schemas"
	"github.com/yhwbach/makerble/internal/utils"
)

// @Summary Register new user
// @Description Register a new doctor or receptionist
// @Tags auth
// @Accept json
// @Produce json
// @Param user body schemas.UserRegister true "User registration info"
// @Success 201 {object} schemas.UserRegisterResponse
// @Failure 400 {object} ErrorResponse
// @Router /register [post]
func (a *Application) registerHandler(w http.ResponseWriter, r *http.Request) {
	var user schemas.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if user email already exists
	existingUserByEmail, err := a.Repo.Users.FindByEmail(r.Context(), user.Email)

	if err != nil && existingUserByEmail != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if existingUserByEmail != nil {
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Check if user username already exists
	existingUserByUsername, err := a.Repo.Users.FindByUsername(r.Context(), user.Username)
	if err != nil && existingUserByUsername != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if existingUserByUsername != nil {
		respondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing registration")
		return
	}

	userID, err := a.Repo.Users.Create(r.Context(), &user, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing request")
		return
	}

	respondWithJSON(w, http.StatusCreated, schemas.UserRegisterResponse{
		UserID:   userIDUUID,
		Message: "User registered successfully",
	})
}

// @Summary Login user
// @Description Login for doctors and receptionists
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body schemas.UserLogin true "Login credentials"
// @Success 200 {object} schemas.TokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (a *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var login schemas.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := a.Repo.Users.FindByUsername(r.Context(), login.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := utils.CheckPassword(login.Password, user.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := a.JWTManager.GenerateToken(user.ID, string(user.UserType))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	respondWithJSON(w, http.StatusOK, schemas.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		UserType:    string(user.UserType),
	})
}

// @Summary Logout user
// @Description Logout current user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Router /logout [post]
func (a *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	tokenStr, err := utils.ExtractBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	// Calculate token expiration time from claims
	var expirationTime time.Time
	if exp, ok := claims["exp"]; ok {
		if expFloat, ok := exp.(float64); ok {
			expirationTime = time.Unix(int64(expFloat), 0)
		}
	}

	if expirationTime.IsZero() {
		expirationTime = time.Now().Add(a.Config.JWT.Expiry)
	}

	// Invalidate the token
	if err := a.Repo.Tokens.InvalidateToken(r.Context(), tokenStr, expirationTime); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error invalidating token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
