package user

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (ur UserRegisterDTO) Validate() error {
	return validation.ValidateStruct(&ur,
		validation.Field(&ur.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&ur.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&ur.FirstName, validation.Required),
		validation.Field(&ur.LastName, validation.Required),
		validation.Field(&ur.LastName, validation.Required),
	)
}

func (c Credentials) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 50)),
	)
}
