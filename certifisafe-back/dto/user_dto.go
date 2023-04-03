package dto

import (
	"certifisafe-back/model"
	"strings"
)

type UserRegisterDTO struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}

type Credentials struct {
	Email    string
	Password string
}

func UserRegisterDTOtoModel(u *UserRegisterDTO) *model.User {
	return &model.User{
		Id:        0,
		Email:     strings.TrimSpace(u.Email),
		Password:  strings.TrimSpace(u.Password),
		FirstName: strings.TrimSpace(u.FirstName),
		LastName:  strings.TrimSpace(u.LastName),
		Phone:     strings.TrimSpace(u.Phone),
		IsAdmin:   false,
	}
}
