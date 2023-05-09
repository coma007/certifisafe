package auth

import (
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AuthController struct {
	service AuthService
}

func NewAuthHandler(cs AuthService) *AuthController {
	return &AuthController{service: cs}
}

func (ah *AuthController) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var credentials user.Credentials
	err := utils.ReadRequestBody(w, r, &credentials)
	if err != nil {
		return
	}

	token, err := ah.service.Login(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, token, http.StatusOK)
}

func (ah *AuthController) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var u user.UserRegisterDTO
	err := utils.ReadRequestBody(w, r, &u)
	if err != nil {
		return
	}

	newUser, err := ah.service.Register(user.UserRegisterDTOtoModel(&u))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, user.ModelToUserBaseDTO(newUser), http.StatusOK)
}

func (ah *AuthController) PasswordRecoveryRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request password_recovery.PasswordRecoveryRequestDTO
	err := utils.ReadRequestBody(w, r, &request)
	if err != nil {
		return
	}

	err = ah.service.RequestPasswordRecoveryToken(request.Email, request.Type)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *AuthController) PasswordRecovery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var request password_recovery.PasswordResetDTO
	err := utils.ReadRequestBody(w, r, &request)
	if err != nil {
		return
	}

	err = ah.service.PasswordRecovery(password_recovery.PasswordResetDTOtoModel(&request))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ah *AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := ps.ByName("verificationCode")
	err := ah.service.VerifyEmail(code)
	if err != nil {
		http.Error(w, "Email verification failed", getAuthErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email successfully verified"))
}

func getAuthErrorStatus(err error) int {
	if errors.Is(err, ErrBadCredentials) ||
		errors.Is(err, ErrTakenEmail) ||
		errors.Is(err, ErrWrongEmailFormat) ||
		errors.Is(err, ErrEmptyName) ||
		errors.Is(err, ErrWrongPhoneFormat) ||
		errors.Is(err, ErrWrongPasswordFormat) ||
		errors.Is(err, ErrCodeUsed) ||
		errors.Is(err, ErrCodeNotFound) ||
		errors.Is(err, ErrNotActivated) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
