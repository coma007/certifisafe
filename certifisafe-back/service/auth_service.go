package service

import (
	"certifisafe-back/model"
	"certifisafe-back/repository"
	"certifisafe-back/utils"
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strings"
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
}

type AuthService struct {
	repository                  repository.IUserRepository
	verificationTokenCharacters string
}

func NewAuthService(repository repository.IUserRepository) *AuthService {
	return &AuthService{repository: repository,
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
	return s.repository.GetUserByEmail(email)
}

func (s *AuthService) Register(user *model.User) (*model.User, error) {
	_, err := s.validateRegistrationData(user)
	if err != nil {
		return &model.User{}, err
	}
	_, err = s.repository.GetUserByEmail(user.Email)
	if err != nil {
		if err == repository.ErrNoUserWithEmail {
			passwordBytes, err := s.hashToken(user.Password)
			utils.CheckError(err)
			user.Password = string(passwordBytes)
			createdUser, err := s.repository.CreateUser(0, *user)
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
