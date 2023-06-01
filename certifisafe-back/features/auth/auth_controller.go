package auth

import (
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Usao")
	var credentials user.Credentials
	err := utils.ReadRequestBody(w, r, &credentials)
	if err != nil {
		return
	}

	token, err := controller.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, token, http.StatusOK)
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {

	var u user.UserRegisterDTO
	err := utils.ReadRequestBody(w, r, &u)
	if err != nil {
		return
	}

	newUser, err := controller.authService.Register(user.UserRegisterDTOtoModel(&u))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, user.ModelToUserBaseDTO(newUser), http.StatusOK)
}

func (controller *AuthController) PasswordRecoveryRequest(w http.ResponseWriter, r *http.Request) {
	var request password_recovery.PasswordRecoveryRequestDTO
	err := utils.ReadRequestBody(w, r, &request)
	if err != nil {
		return
	}

	err = controller.authService.RequestPasswordRecoveryToken(request.Email, request.Type)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *AuthController) PasswordRecovery(w http.ResponseWriter, r *http.Request) {
	var request password_recovery.PasswordResetDTO
	err := utils.ReadRequestBody(w, r, &request)
	if err != nil {
		return
	}

	err = controller.authService.PasswordRecovery(password_recovery.PasswordResetDTOtoModel(&request))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)

	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (controller *AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	code := utils.ReadVerificationCodeFromUrl(w, r)
	err := controller.authService.VerifyEmail(code)
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
		errors.Is(err, ErrNotActivated) ||
		errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrPasswordChange) {
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
