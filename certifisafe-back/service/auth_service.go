package service

import (
	"certifisafe-back/model"
	"certifisafe-back/repository"
	"certifisafe-back/utils"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var (
	ErrBadCredentials = errors.New("bad username or password")
	ErrTakenEmail     = errors.New("email already taken")
)

type IAuthService interface {
	Login(email string, password string) (string, error)
	ValidateToken(tokenString string) (bool, error)
	Register(user model.User) (model.User, error)
}

type AuthService struct {
	repository repository.IUserRepository
}

func NewAuthService(repository repository.IUserRepository) *AuthService {
	return &AuthService{repository: repository}
}

var jwtKey = []byte("secret-key")

type Claims struct {
	Email string
	jwt.StandardClaims
}

func (s *AuthService) Login(email string, password string) (string, error) {
	user, err := s.repository.GetUserByEmail(email)
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

func (s *AuthService) Register(user model.User) (model.User, error) {

	_, err := s.repository.GetUserByEmail(user.Email)
	if err != nil {
		if err == repository.ErrNoUserWithEmail {
			passwordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
			utils.CheckError(err)
			user.Password = string(passwordBytes)
			createdUser, err := s.repository.CreateUser(0, user)
			if err != nil {
				return model.User{}, err
			}
			return createdUser, nil
		} else {
			return model.User{}, ErrTakenEmail
		}
	}

	return model.User{}, ErrTakenEmail
}

func (s *AuthService) ValidateToken(tokenString string) (bool, error) {

	tokens := strings.Split(tokenString, " ")
	tokenString, schema := tokens[0], tokens[1]
	//tokenString = tokenString[1 : len(tokenString)-1]

	if schema != "Bearer" {
		return false, errors.New("New")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}
