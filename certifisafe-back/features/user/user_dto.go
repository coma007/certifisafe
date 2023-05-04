package user

import (
	"strings"
)

type UserBaseDTO struct {
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

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

func UserRegisterDTOtoModel(u *UserRegisterDTO) *User {
	return &User{
		Email:     strings.TrimSpace(u.Email),
		Password:  strings.TrimSpace(u.Password),
		FirstName: strings.TrimSpace(u.FirstName),
		LastName:  strings.TrimSpace(u.LastName),
		Phone:     strings.TrimSpace(u.Phone),
		IsAdmin:   false,
	}
}

func ModelToUserBaseDTO(u *User) *UserBaseDTO {
	return &UserBaseDTO{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
	}
}
