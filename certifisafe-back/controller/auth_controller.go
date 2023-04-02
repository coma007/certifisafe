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
		http.Error(w, err.Error(), getErrorLoginStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)

	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func getErrorLoginStatus(err error) int {
	if errors.Is(err, service.ErrBadCredentials) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func (ah *AuthHandler) Validate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := r.Header.Get("Authorization")

	ah.service.ValidateToken(token)
}
