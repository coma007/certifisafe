package auth

import (
	"bytes"
	password_recovery2 "certifisafe-back/features/password_recovery"
	user2 "certifisafe-back/features/user"
	"certifisafe-back/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/smtp"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unicode"
)

var (
	ErrNoUserWithEmail     = errors.New("no user for given email")
	ErrBadCredentials      = errors.New("bad username or password")
	ErrNotActivated        = errors.New("account is not activated")
	ErrTakenEmail          = errors.New("email already taken")
	ErrWrongEmailFormat    = errors.New("not valid email")
	ErrEmptyName           = errors.New("name cannot be empty")
	ErrWrongPhoneFormat    = errors.New("not valid phone")
	ErrWrongPasswordFormat = errors.New("not valid password")
	ErrCodeUsed            = errors.New("verification code is used")
	ErrCodeNotFound        = errors.New("verification code cannot be found")
)

type AuthService interface {
	Login(email string, password string) (string, error)
	ValidateToken(tokenString string) (bool, error)
	Register(user *user2.User) (*user2.User, error)
	VerifyEmail(verificationCode string) error
	GetUserFromToken(tokenString string) user2.User
	GetClaims(tokenString string) (*jwt.Token, *Claims, bool, error)
	GetUserByEmail(email string) (user2.User, error)
	RequestPasswordRecoveryToken(email string) error
	PasswordRecovery(request *password_recovery2.PasswordRecovery) error
}

// TODO check if everything works - Duti (Bobi made changes)
// TODO separate mails to email_service.go
type DefaultAuthService struct {
	userRepository              user2.UserRepository
	passwordRecoveryRepository  password_recovery2.PasswordRecoveryRepository
	verificationRepository      VerificationRepository
	verificationTokenCharacters string
}

func NewDefaultAuthService(userRepository user2.UserRepository, passwordRecoveryRepository password_recovery2.PasswordRecoveryRepository,
	verificationRepository VerificationRepository) *DefaultAuthService {
	return &DefaultAuthService{userRepository: userRepository,
		passwordRecoveryRepository:  passwordRecoveryRepository,
		verificationRepository:      verificationRepository,
		verificationTokenCharacters: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"}
}

var jwtKey = []byte("secret-key")

type Claims struct {
	Email string
	jwt.StandardClaims
}

func (s *DefaultAuthService) Login(email string, password string) (string, error) {
	user, err := s.GetUserByEmail(email)
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

func (s *DefaultAuthService) GetUserByEmail(email string) (user2.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *DefaultAuthService) Register(u *user2.User) (*user2.User, error) {
	u.IsActive = false
	_, err := s.validateRegistrationData(u)
	if err != nil {
		return &user2.User{}, err
	}
	_, err = s.userRepository.GetUserByEmail(u.Email)
	if err != nil {
		if err == ErrNoUserWithEmail {
			passwordBytes, err := s.hashToken(u.Password)
			utils.CheckError(err)
			u.Password = string(passwordBytes)
			createdUser, err := s.userRepository.CreateUser(*u)
			if err != nil {
				return &user2.User{}, err
			}

			//add phone option
			s.sendVerification(u)

			return &createdUser, nil
		} else {
			return &user2.User{}, ErrTakenEmail
		}
	}

	return &user2.User{}, ErrTakenEmail
}

func (s *DefaultAuthService) ValidateToken(tokenString string) (bool, error) {

	token, _, b, err2 := s.GetClaims(tokenString)
	if err2 != nil {
		return b, err2
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}

func (s *DefaultAuthService) GetUserFromToken(tokenString string) user2.User {
	_, claims, _, _ := s.GetClaims(tokenString)
	email := claims.Email
	user, _ := s.GetUserByEmail(email)
	return user
}

func (s *DefaultAuthService) GetClaims(tokenString string) (*jwt.Token, *Claims, bool, error) {
	tokens := strings.Split(tokenString, " ")
	schema, tokenString := tokens[0], tokens[1]
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

func (s *DefaultAuthService) RequestPasswordRecoveryToken(email string) error {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return err
	}

	to := []string{user.Email}

	templateFile, _ := filepath.Abs("utils/passwordRecovery.html")
	t, err := template.ParseFiles(templateFile)

	if err != nil {
		return err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Password recovery \n%s\n\n", mimeHeaders)))

	verificationToken, err := s.getVerificationToken(4, false)

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

	_, err = s.passwordRecoveryRepository.CreateRequest(1, password_recovery2.PasswordRecoveryRequest{Email: user.Email, Code: string(verificationToken)})
	if err != nil {
		return err
	}

	err = s.sendMail(to, body)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultAuthService) PasswordRecovery(request *password_recovery2.PasswordRecovery) error {
	//token, err := s.hashToken(request.Code)
	//if err != nil {
	//	return err
	//}
	r, err := s.passwordRecoveryRepository.GetRequestByCode(string(request.Code))
	if err != nil {
		return err
	}
	if r.IsUsed {
		return ErrCodeUsed
	}

	user, err := s.userRepository.GetUserByEmail(r.Email)

	//verify password
	if !s.verifyPassword(request.NewPassword) {
		return ErrWrongPasswordFormat
	}

	hashedPassword, err := s.hashToken(request.NewPassword)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	s.userRepository.UpdateUser(int32(user.ID), user)
	s.passwordRecoveryRepository.UseRequestsForEmail(user.Email)
	return nil
}

func (s *DefaultAuthService) VerifyEmail(verificationCode string) error {
	verification, err := s.verificationRepository.GetVerificationByCode(verificationCode)
	if err != nil {
		return ErrCodeNotFound
	}
	user, err := s.userRepository.GetUserByEmail(verification.Email)
	if err != nil {
		return err
	}
	user.IsActive = true
	_, err = s.userRepository.UpdateUser(int32(user.ID), user)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultAuthService) sendMail(to []string, body bytes.Buffer) error {
	from := "ftn.project.usertest@gmail.com"
	password := "zmiwmhfweojejlqy"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	go smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	return nil
}

func (s *DefaultAuthService) sendVerification(user *user2.User) error {
	to := []string{user.Email}

	templateFile, _ := filepath.Abs("utils/emailVerification.html")
	t, err := template.ParseFiles(templateFile)

	if err != nil {
		return err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Email verification \n%s\n\n", mimeHeaders)))

	verificationToken, err := s.getVerificationToken(10, true)
	if err != nil {
		return err
	}

	_, err = s.verificationRepository.CreateVerification(0, Verification{
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

	s.sendMail(to, body)
	return nil
}

func (s *DefaultAuthService) getVerificationToken(length int, verification bool) (string, error) {

	verificationString := ""
	for true {
		for i := 0; i < length; i++ {
			nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(s.verificationTokenCharacters))))
			if err != nil {
				return "", err
			}
			verificationString += string(s.verificationTokenCharacters[nBig.Int64()])
		}
		var err error
		if verification {
			_, err = s.verificationRepository.GetVerificationByCode(verificationString)
		} else {
			_, err = s.passwordRecoveryRepository.GetRequestByCode(verificationString)
		}
		if err != nil {
			break
		} else {
			verificationString = ""
		}
	}
	return verificationString, nil
}

func (s *DefaultAuthService) hashToken(password string) ([]byte, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return passwordBytes, err
}

func (s *DefaultAuthService) validateRegistrationData(u *user2.User) (bool, error) {
	match, err := regexp.Match("^[\\w-\\+\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(u.Email))
	if err != nil || !match || u.Email == "" {
		return false, ErrWrongEmailFormat
	} else if u.FirstName == "" || u.LastName == "" {
		return false, ErrEmptyName
	}

	if !s.verifyPassword(u.Password) {
		return false, ErrWrongPasswordFormat
	}
	match, err = regexp.Match("^[0-9]*$", []byte(u.Phone))
	if err != nil || !match || u.Phone == "" || (len(u.Phone) != 9 && len(u.Phone) != 10) {
		return false, ErrWrongPhoneFormat
	}
	return true, nil
}

func (s *DefaultAuthService) verifyPassword(password string) bool {
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
