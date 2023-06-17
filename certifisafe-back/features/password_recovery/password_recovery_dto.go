package password_recovery

type PasswordRecoveryRequestDTO struct {
	Email string
	Type  int
}

type PasswordResetDTO struct {
	VerificationCode string
	NewPassword      string
}

func PasswordResetDTOtoModel(dto *PasswordResetDTO) *PasswordRecovery {
	return &PasswordRecovery{
		Code:        dto.VerificationCode,
		NewPassword: dto.NewPassword,
	}
}
