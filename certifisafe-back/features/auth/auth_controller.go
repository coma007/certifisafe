package auth

import (
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type AuthController struct {
	authService AuthService
	oauthConfig *oauth2.Config
}

func NewAuthController(authService AuthService) *AuthController {
	config := utils.Config()
	return &AuthController{
		authService: authService,
		oauthConfig: &oauth2.Config{
			ClientID:     config["oauth-client-id"],
			ClientSecret: config["oauth-client-secret"],
			RedirectURL:  config["oauth-redirect-url"],
			Scopes:       []string{"openid", "email", "profile"}, // Add required scopes
			Endpoint:     google.Endpoint,
		},
	}
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials user.Credentials
	err := utils.ReadRequestBody(w, r, &credentials)
	if err != nil {
		return
	}

	// template for validation
	err = credentials.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := controller.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, token, http.StatusNoContent)
}

func (controller *AuthController) TwoFactorAuth(w http.ResponseWriter, r *http.Request) {
	var code CodeDTO
	err := utils.ReadRequestBody(w, r, &code)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	token, err := controller.authService.TwoFactorAuth(code.VerificationCode)
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, token, http.StatusOK)
}

func (controller *AuthController) OauthLogin(w http.ResponseWriter, r *http.Request) {
	authURL := controller.oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {

	var u user.UserRegisterDTO
	err := utils.ReadRequestBody(w, r, &u)
	if err != nil {
		return
	}

	err = u.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := controller.authService.Register(user.UserRegisterDTOtoModel(&u))
	if err != nil {
		http.Error(w, err.Error(), getAuthErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, user.ModelToUserBaseDTO(newUser), http.StatusOK)
}

func (controller *AuthController) OauthRegister(writer http.ResponseWriter, request *http.Request) {

}

func (controller *AuthController) OauthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	token, err := controller.oauthConfig.Exchange(context.Background(), state)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// You can use the token to retrieve user information or perform any necessary actions

	// Example: Retrieve user information
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Provider string `json:"provider"`
	}

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, "Failed to parse user information", http.StatusInternalServerError)
		return
	}

	// Use the user information as needed
	fmt.Println("Email:", userInfo.Email)
	fmt.Println("Name:", userInfo.Name)

	// Redirect the user to the desired page or perform further actions
	http.Redirect(w, r, controller.oauthConfig.RedirectURL, http.StatusTemporaryRedirect)
}

func (controller *AuthController) PasswordRecoveryRequest(w http.ResponseWriter, r *http.Request) {
	var request password_recovery.PasswordRecoveryRequestDTO
	err := utils.ReadRequestBody(w, r, &request)
	if err != nil {
		return
	}

	err = request.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.authService.RequestPasswordRecoveryToken(request.Email, request.Type, 0)
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

	err = request.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	if strings.Compare("", strings.TrimSpace(code)) == 0 {
		http.Error(w, "Code cannot be empty string", http.StatusBadRequest)
		return
	}
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
