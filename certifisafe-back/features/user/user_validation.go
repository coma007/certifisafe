package user

import (
	"certifisafe-back/features/password_recovery"
	_ "certifisafe-back/features/password_recovery"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (ur UserRegisterDTO) Validate() error {
	return validation.ValidateStruct(&ur,
		validation.Field(&ur.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&ur.Password, validation.Required, validation.Length(8, 50),
			password_recovery.PasswordValidation[0], password_recovery.PasswordValidation[1], password_recovery.PasswordValidation[2]),
		validation.Field(&ur.FirstName, validation.Required),
		validation.Field(&ur.LastName, validation.Required),
		validation.Field(&ur.LastName, validation.Required),
	)
}

func (c Credentials) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 50),
			password_recovery.PasswordValidation[0], password_recovery.PasswordValidation[1], password_recovery.PasswordValidation[2]),
	)
}
