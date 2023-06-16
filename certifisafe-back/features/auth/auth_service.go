package auth

import (
	"bytes"
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/big"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unicode"
)

var (
	ErrNoUserWithEmail      = errors.New("no user for given email")
	ErrBadCredentials       = errors.New("bad username or password")
	ErrNotActivated         = errors.New("account is not activated")
	ErrPasswordChange       = errors.New("password needs to be changed, email has been sent")
	ErrPasswordNotAvailable = errors.New("cannot use old password")
	ErrTakenEmail           = errors.New("email already taken")
	ErrWrongEmailFormat     = errors.New("not valid email")
	ErrEmptyName            = errors.New("name cannot be empty")
	ErrWrongPhoneFormat     = errors.New("not valid phone")
	ErrWrongPasswordFormat  = errors.New("not valid password")
	ErrCodeUsed             = errors.New("verification code is used")
	ErrCodeNotFound         = errors.New("verification code cannot be found")
)

type AuthService interface {
	Login(email string, password string) (string, error)
	ValidateToken(tokenString string) (bool, error)
	Register(user *user.User) (*user.User, error)
	VerifyEmail(verificationCode string) error
	GetUserFromToken(tokenString string) user.User
	GetClaims(tokenString string) (*jwt.Token, *Claims, bool, error)
	GetUserByEmail(email string) (user.User, error)
	RequestPasswordRecoveryToken(email string, t int) error
	PasswordRecovery(request *password_recovery.PasswordRecovery) error
}

type DefaultAuthService struct {
	mailService                 MailService
	userRepository              user.UserRepository
	passwordRecoveryRepository  password_recovery.PasswordRecoveryRepository
	passwordHistoryRepository   password_recovery.PasswordHistoryRepository
	verificationRepository      VerificationRepository
	verificationTokenCharacters string
}

func NewDefaultAuthService(mailService MailService, userRepo user.UserRepository, passwordRecoveryRepo password_recovery.PasswordRecoveryRepository,
	passwordHistoryRepo password_recovery.PasswordHistoryRepository,
	verificationRepo VerificationRepository) *DefaultAuthService {
	return &DefaultAuthService{mailService: mailService,
		userRepository:              userRepo,
		passwordRecoveryRepository:  passwordRecoveryRepo,
		verificationRepository:      verificationRepo,
		passwordHistoryRepository:   passwordHistoryRepo,
		verificationTokenCharacters: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"}
}

var jwtKey = []byte("secret-key")

type Claims struct {
	Email string
	jwt.StandardClaims
}

func (service *DefaultAuthService) Login(email string, password string) (string, error) {
	user, err := service.GetUserByEmail(email)
	if err != nil {
		if err == ErrNoUserWithEmail {
			return "", ErrBadCredentials
		} else {
			return "", err
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		if !user.IsActive {
			return "", ErrNotActivated
		}
		if time.Since(user.LastPasswordSet).Milliseconds() > 1000*60*60 {
			err := service.RequestPasswordRecoveryToken(user.Email, 0)
			if err != nil {
				return "", err
			}
			return "", ErrPasswordChange
		}
		expirationTime := time.Now().Add(time.Minute * 60)

		claims := &Claims{
			Email:          email,
			StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)

		utils.CheckError(err)

		return tokenString, nil
	}

	return "", ErrBadCredentials
}

func (service *DefaultAuthService) GetUserByEmail(email string) (user.User, error) {
	return service.userRepository.GetUserByEmail(email)
}

func (service *DefaultAuthService) Register(u *user.User) (*user.User, error) {
	u.IsActive = false
	_, err := service.validateRegistrationData(u)
	if err != nil {
		return &user.User{}, err
	}
	_, err = service.userRepository.GetUserByEmail(u.Email)
	if err == gorm.ErrRecordNotFound {
		passwordBytes, err := service.hashToken(u.Password)
		utils.CheckError(err)
		u.Password = string(passwordBytes)
		createdUser, err := service.userRepository.CreateUser(*u)
		if err != nil {
			return &user.User{}, err
		}

		//add phone option
		service.sendVerification(u)

		return &createdUser, nil
	}

	return &user.User{}, ErrTakenEmail
}

func (service *DefaultAuthService) ValidateToken(tokenString string) (bool, error) {

	token, _, b, err2 := service.GetClaims(tokenString)
	if err2 != nil {
		return b, err2
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}

func (service *DefaultAuthService) GetUserFromToken(tokenString string) user.User {
	_, claims, _, _ := service.GetClaims(tokenString)
	email := claims.Email
	user, _ := service.GetUserByEmail(email)
	return user
}

func (service *DefaultAuthService) GetClaims(tokenString string) (*jwt.Token, *Claims, bool, error) {
	tokens := strings.Split(tokenString, " ")
	schema, tokenString := tokens[0], tokens[1]
	if strings.ToLower(strings.TrimSpace(tokenString)) == `bearer` {
		tokenString, schema = schema, tokenString
	}
	//tokenString = tokenString[1 : len(tokenString)-1]

	if schema != "Bearer" {
		return nil, nil, false, errors.New("New")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		return nil, claims, false, err
	}
	return token, claims, false, nil
}

// 0 - email, 1 - phone
func (service *DefaultAuthService) RequestPasswordRecoveryToken(value string, t int) error {
	var user user.User
	var err error
	if t == 0 {
		user, err = service.userRepository.GetUserByEmail(value)
	} else {
		user, err = service.userRepository.GetUserByPhone(value)
	}
	if err != nil {
		return err
	}

	to := []string{user.Email}

	templateFile, _ := filepath.Abs("resources/templates/passwordRecovery.html")
	temp, err := template.ParseFiles(templateFile)

	if err != nil {
		return err
	}

	var body bytes.Buffer
	//mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//body.Write([]byte(fmt.Sprintf("Subject: Password recovery \n%service\n\n", mimeHeaders)))

	verificationToken, err := service.getVerificationToken(4, false)

	if err != nil {
		return err
	}

	temp.Execute(&body, struct {
		Name string
		Code string
	}{
		Name: user.FirstName + " " + user.LastName,
		Code: verificationToken,
	})

	_, err = service.passwordRecoveryRepository.CreateRequest(1, password_recovery.PasswordRecoveryRequest{Email: user.Email, Code: string(verificationToken)})
	if err != nil {
		return err
	}

	if t == 1 {
		err = service.sendSMS(verificationToken)
		if err != nil {
			return err
		}
		return nil
	}

	err = service.mailService.SendMail(to, body)
	if err != nil {
		return err
	}
	return nil
}

func (service *DefaultAuthService) sendSMS(verificationToken string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TwillioApiUsername,
		Password: TwillioApiPassword,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo("+38162711935")
	params.SetFrom("+12708136240")
	params.SetBody("Here is your one time recovery code: " + verificationToken)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}
	return err
}

func (service *DefaultAuthService) PasswordRecovery(request *password_recovery.PasswordRecovery) error {
	//token, err := service.hashToken(request.Code)
	//if err != nil {
	//	return err
	//}
	r, err := service.passwordRecoveryRepository.GetRequestByCode(string(request.Code))
	if err != nil {
		return err
	}
	if r.IsUsed {
		return ErrCodeUsed
	}

	user, err := service.userRepository.GetUserByEmail(r.Email)

	//verify password
	if !service.verifyPassword(request.NewPassword) {
		return ErrWrongPasswordFormat
	}

	if service.isPasswordUsed(user.Email, request.NewPassword) {
		return ErrPasswordNotAvailable
	}
	hashedPassword, err := service.hashToken(request.NewPassword)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.LastPasswordSet = time.Now()
	service.userRepository.UpdateUser(user.ID, user)
	service.passwordRecoveryRepository.UseRequestsForEmail(user.Email)
	return nil
}

func (service *DefaultAuthService) VerifyEmail(verificationCode string) error {
	verification, err := service.verificationRepository.GetVerificationByCode(verificationCode)
	if err != nil {
		return ErrCodeNotFound
	}
	user, err := service.userRepository.GetUserByEmail(verification.Email)
	if err != nil {
		return err
	}
	user.IsActive = true
	_, err = service.userRepository.UpdateUser(user.ID, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *DefaultAuthService) sendVerification(user *user.User) error {
	to := []string{user.Email}

	templateFile, _ := filepath.Abs("resources/templates/emailVerification.html")
	t, err := template.ParseFiles(templateFile)

	if err != nil {
		return err
	}

	var body bytes.Buffer
	//mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//body.Write([]byte(fmt.Sprintf("Subject: Email verification \n%service\n\n", mimeHeaders)))

	verificationToken, err := service.getVerificationToken(10, true)
	if err != nil {
		return err
	}

	_, err = service.verificationRepository.CreateVerification(0, Verification{
		Email: user.Email,
		Code:  verificationToken,
	})

	if err != nil {
		return err
	}

	t.Execute(&body, struct {
		Name string
		Code string
	}{
		Name: user.FirstName + " " + user.LastName,
		Code: verificationToken,
	})

	service.mailService.SendMail(to, body)
	return nil
}

func (service *DefaultAuthService) getVerificationToken(length int, verification bool) (string, error) {

	verificationString := ""
	for true {
		for i := 0; i < length; i++ {
			nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(service.verificationTokenCharacters))))
			if err != nil {
				return "", err
			}
			verificationString += string(service.verificationTokenCharacters[nBig.Int64()])
		}
		var err error
		if verification {
			_, err = service.verificationRepository.GetVerificationByCode(verificationString)
		} else {
			_, err = service.passwordRecoveryRepository.GetRequestByCode(verificationString)
		}
		if err != nil {
			break
		} else {
			verificationString = ""
		}
	}
	return verificationString, nil
}

func (service *DefaultAuthService) hashToken(password string) ([]byte, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return passwordBytes, err
}

func (service *DefaultAuthService) validateRegistrationData(u *user.User) (bool, error) {
	match, err := regexp.Match("^[\\w-\\+\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(u.Email))
	if err != nil || !match || u.Email == "" {
		return false, ErrWrongEmailFormat
	} else if u.FirstName == "" || u.LastName == "" {
		return false, ErrEmptyName
	}

	if !service.verifyPassword(u.Password) {
		return false, ErrWrongPasswordFormat
	}
	match, err = regexp.Match("^[0-9]*$", []byte(u.Phone))
	if err != nil || !match || u.Phone == "" || (len(u.Phone) != 9 && len(u.Phone) != 10) {
		return false, ErrWrongPhoneFormat
	}
	return true, nil
}

func (service *DefaultAuthService) verifyPassword(password string) bool {
	number, upper, lower := false, false, false

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		default:
			continue
		}
	}
	return number && upper && lower && len(password) >= 8
}

func (service *DefaultAuthService) isPasswordUsed(email string, password string) bool {

	history, err := service.passwordHistoryRepository.GetHistoryByEmail(email)
	if err != nil {
		return true
	}

	for _, element := range history {
		if bcrypt.CompareHashAndPassword([]byte(element.ForbiddenPassword), []byte(password)) == nil {
			return true
		}
	}

	numberOfPasswords := len(history)
	if numberOfPasswords == 2 {
		firstPassword := history[0]
		for _, element := range history {
			if element.ID < firstPassword.ID {
				firstPassword = element
			}
		}
		service.passwordHistoryRepository.DeleteHistory(int32(firstPassword.ID))
	}
	hashedPassword, err := service.hashToken(password)
	if err != nil {
		return true
	}
	_, err = service.passwordHistoryRepository.CreateHistory(0, password_recovery.PasswordHistory{
		Model:             gorm.Model{},
		Deleted:           gorm.DeletedAt{},
		UserEmail:         email,
		ForbiddenPassword: string(hashedPassword),
	})
	if err != nil {
		return true
	}
	return false
}
