package dto

type PasswordRecoveryRequestDTO struct {
	Email string
}

type PasswordResetDTO struct {
	VerificationCode string
	NewPassword      string
}
