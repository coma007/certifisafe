package controller

import (
	"certifisafe-back/dto"
	"certifisafe-back/service"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AuthHandler struct {
	service service.IAuthService
}

func NewAuthHandler(cs service.IAuthService) *AuthHandler {
	return &AuthHandler{service: cs}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var credentials dto.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	token, err := ah.service.Login(credentials.Email, credentials.Password)

	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)

	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user dto.UserRegisterDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}
	newUser, err := ah.service.Register(dto.UserRegisterDTOtoModel(&user))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dto.ModelToUserBaseDTO(newUser))

	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (ah *AuthHandler) PasswordRecoveryRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request dto.PasswordRecoveryRequestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	err = ah.service.RequestPasswordRecoveryToken(request.Email)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ah *AuthHandler) PasswordRecovery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request dto.PasswordResetDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	err = ah.service.PasswordRecovery(dto.PasswordResetDTOtoModel(&request))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getAuthErrorStatus(err error) int {
	if errors.Is(err, service.ErrBadCredentials) ||
		errors.Is(err, service.ErrTakenEmail) ||
		errors.Is(err, service.ErrWrongEmailFormat) ||
		errors.Is(err, service.ErrEmptyName) ||
		errors.Is(err, service.ErrWrongPhoneFormat) ||
		errors.Is(err, service.ErrWrongPasswordFormat) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
