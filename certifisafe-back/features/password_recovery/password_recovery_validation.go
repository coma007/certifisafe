package password_recovery

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (pr PasswordRecoveryRequestDTO) Validate() error {
	return validation.ValidateStruct(&pr,
		validation.Field(&pr.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&pr.Type, validation.Required, is.Int),
	)
}

func (pr PasswordResetDTO) Validate() error {
	return validation.ValidateStruct(&pr,
		validation.Field(&pr.NewPassword, validation.Required, validation.Length(8, 50)),
		validation.Field(&pr.VerificationCode, validation.Required, is.Alphanumeric),
	)
}
