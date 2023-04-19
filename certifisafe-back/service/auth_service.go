package service

import (
	"bytes"
	"certifisafe-back/model"
	"certifisafe-back/repository"
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
	ErrBadCredentials      = errors.New("bad username or password")
	ErrTakenEmail          = errors.New("email already taken")
	ErrWrongEmailFormat    = errors.New("not valid email")
	ErrEmptyName           = errors.New("name cannot be empty")
	ErrWrongPhoneFormat    = errors.New("not valid phone")
	ErrWrongPasswordFormat = errors.New("not valid password")
)

type IAuthService interface {
	Login(email string, password string) (string, error)
	ValidateToken(tokenString string) (bool, error)
	Register(user *model.User) (*model.User, error)
	GetClaims(tokenString string) (*jwt.Token, *Claims, bool, error)
	GetUserByEmail(email string) (model.User, error)
	RequestPasswordRecoveryToken(email string) error
}

type AuthService struct {
	userRepository              repository.IUserRepository
	passwordRecoveryRepository  repository.IPasswordRecoveryRepository
	verificationTokenCharacters string
}

func NewAuthService(userRepository repository.IUserRepository, passwordRecoveryRepository repository.IPasswordRecoveryRepository) *AuthService {
	return &AuthService{userRepository: userRepository,
		passwordRecoveryRepository:  passwordRecoveryRepository,
		verificationTokenCharacters: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"}
}

var jwtKey = []byte("secret-key")

type Claims struct {
	Email string
	jwt.StandardClaims
}

func (s *AuthService) Login(email string, password string) (string, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		if err == repository.ErrNoUserWithEmail {
			return "", ErrBadCredentials
		} else {
			return "", err
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
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

func (s *AuthService) GetUserByEmail(email string) (model.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *AuthService) Register(user *model.User) (*model.User, error) {
	_, err := s.validateRegistrationData(user)
	if err != nil {
		return &model.User{}, err
	}
	_, err = s.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		if err == repository.ErrNoUserWithEmail {
			passwordBytes, err := s.hashToken(user.Password)
			utils.CheckError(err)
			user.Password = string(passwordBytes)
			createdUser, err := s.userRepository.CreateUser(0, *user)
			if err != nil {
				return &model.User{}, err
			}
			return &createdUser, nil
		} else {
			return &model.User{}, ErrTakenEmail
		}
	}

	return &model.User{}, ErrTakenEmail
}

func (s *AuthService) ValidateToken(tokenString string) (bool, error) {

	token, _, b, err2 := s.GetClaims(tokenString)
	if err2 != nil {
		return b, err2
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}

func (s *AuthService) GetClaims(tokenString string) (*jwt.Token, *Claims, bool, error) {
	tokens := strings.Split(tokenString, " ")
	tokenString, schema := tokens[0], tokens[1]
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

func (s *AuthService) RequestPasswordRecoveryToken(email string) error {
	from := "ftn.project.usertest@gmail.com"
	password := "zmiwmhfweojejlqy"

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return err
	}

	//request, err := s.passwordRecoveryRepository.GetRequestsByEmail(email)

	//if err == nil {
	//fmt.Println(request)
	//err := s.passwordRecoveryRepository.DeleteRequest(int32(request.Id))
	//if err != nil {
	//	return err
	//}
	//}

	to := []string{user.Email}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	templateFile, _ := filepath.Abs("utils/passwordRecovery.html")
	t, err := template.ParseFiles(templateFile)

	if err != nil {
		return err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Password recovery \n%s\n\n", mimeHeaders)))

	verificationToken, err := s.getVerificationToken()

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

	token, err := s.hashToken(verificationToken)
	if err != nil {
		return err
	}
	_, err = s.passwordRecoveryRepository.CreateRequest(1, model.PasswordRecoveryRequest{Id: 1, Email: user.Email, Code: string(token)})
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	//req, err := s.passwordRecoveryRepository.GetRequestsByEmail(user.Email)
	//if err != nil {
	//	return err
	//}
	//
	//for i := 0; i < len(req); i++ {
	//	fmt.Println(*req[i])
	//}
	//fmt.Println("Email Sent!")

	return nil
}

func (s *AuthService) getVerificationToken() (string, error) {

	verificationString := ""
	for i := 0; i < 4; i++ {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(s.verificationTokenCharacters))))
		if err != nil {
			return "", err
		}
		verificationString += string(s.verificationTokenCharacters[nBig.Int64()])
	}
	return verificationString, nil
}

func (s *AuthService) hashToken(password string) ([]byte, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return passwordBytes, err
}

func (s *AuthService) validateRegistrationData(u *model.User) (bool, error) {
	match, err := regexp.Match("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(u.Email))
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

func (s *AuthService) verifyPassword(password string) bool {
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
