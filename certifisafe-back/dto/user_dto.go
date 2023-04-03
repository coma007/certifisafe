package dto

import "certifisafe-back/model"

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
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		IsAdmin:   false,
	}
}
