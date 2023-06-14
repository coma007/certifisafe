package password_recovery

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
)

func (pr PasswordRecoveryRequestDTO) Validate() error {
	return validation.ValidateStruct(&pr,
		validation.Field(&pr.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&pr.Type, validation.Required, is.Int),
	)
}

func (pr PasswordResetDTO) Validate() error {
	return validation.ValidateStruct(&pr,
		//Minimum eight characters, at least one uppercase letter, one lowercase letter and one number:
		validation.Field(&pr.NewPassword, validation.Required, validation.Length(8, 50),
			PasswordValidation[0], PasswordValidation[1], PasswordValidation[2]),
		validation.Field(&pr.VerificationCode, validation.Required, is.Alphanumeric),
	)
}

var PasswordValidation = [3]*validation.MatchRule{
	validation.Match(regexp.MustCompile("[a-z]+")).Error("must have at least a lowercase letter"),
	validation.Match(regexp.MustCompile("[A-Z]+")).Error("must have at least an uppercase letter"),
	validation.Match(regexp.MustCompile("[0-9]+")).Error("must have at least one number"),
}
