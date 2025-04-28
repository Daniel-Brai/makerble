package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yhwbach/makerble/internal/schemas"
	"github.com/yhwbach/makerble/internal/utils"
)

// @Summary Create patient
// @Description Create a new patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param patient body schemas.PatientCreate true "Patient information"
// @Success 201 {object} schemas.PatientCreateResponse
// @Failure 400,403 {object} ErrorResponse
// @Router /patients [post]
func (a *Application) createPatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient schemas.PatientCreate
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Parse date of birth
	dateOfBirth, err := time.Parse("2006-01-02", patient.DateOfBirth)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing request")
		return
	}

	registeredByID, err := uuid.Parse(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing request")
		return
	}

	// Create a modified patient with the parsed date
	patientToCreate := patient
	patientID, err := a.Repo.Patients.Create(r.Context(), registeredByID, &patientToCreate, dateOfBirth)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating patient")
		return
	}

	respondWithJSON(w, http.StatusCreated, schemas.PatientCreateResponse{
		Message:   "Patient created successfully",
		PatientID: patientID,
	})
}

// @Summary List patients
// @Description Get all patients
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} schemas.PatientListResponse
// @Failure 403,500 {object} ErrorResponse
// @Router /patients [get]
func (a *Application) listPatientsHandler(w http.ResponseWriter, r *http.Request) {
	patients, err := a.Repo.Patients.FindAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching patients")
		return
	}

	respondWithJSON(w, http.StatusOK, schemas.PatientListResponse{
		Patients: patients,
		Total:    len(patients),
		Page:     1,
		PageSize: 10,
	})
}

// @Summary Get patient
// @Description Get patient by ID
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Success 200 {object} models.Patient
// @Failure 403,404,500 {object} ErrorResponse
// @Router /patients/{id} [get]
func (a *Application) getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := a.Repo.Patients.FindByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching patient")
		return
	}
	if patient == nil {
		respondWithError(w, http.StatusNotFound, "Patient not found")
		return
	}

	respondWithJSON(w, http.StatusOK, patient)
}

// @Summary Update patient
// @Description Update patient information (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Param patient body schemas.PatientUpdate true "Patient update information"
// @Success 200 {object} models.Patient
// @Failure 400,403,404,500 {object} ErrorResponse
// @Router /patients/{id} [put]
func (a *Application) updatePatientLimitedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	var update schemas.PatientUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	limitedUpdate := schemas.PatientUpdate{
		FullName:       update.FullName,
		Email:         update.Email,
		Phone:         update.Phone,
		Address:       update.Address,
	}

	patient, err := a.Repo.Patients.UpdateByID(r.Context(), id, &limitedUpdate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating patient")
		return
	}
	if patient == nil {
		respondWithError(w, http.StatusNotFound, "Patient not found")
		return
	}

	respondWithJSON(w, http.StatusOK, patient)
}

// @Summary Update patient medical info
// @Description Update patient medical information (Doctor only)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Param medical_info body schemas.PatientUpdate true "Medical information update"
// @Success 200 {object} models.Patient
// @Failure 400,403,404,500 {object} ErrorResponse
// @Router /patients/{id} [patch]
func (a *Application) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	var update schemas.PatientUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	patient, err := a.Repo.Patients.UpdateByID(r.Context(), id, &update)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating patient medical info")
		return
	}
	if patient == nil {
		respondWithError(w, http.StatusNotFound, "Patient not found")
		return
	}

	respondWithJSON(w, http.StatusOK, patient)
}

// @Summary Delete patient
// @Description Delete a patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Success 204 "No Content"
// @Failure 403,404,500 {object} ErrorResponse
// @Router /patients/{id} [delete]
func (a *Application) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	err = a.Repo.Patients.DeleteByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting patient")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
