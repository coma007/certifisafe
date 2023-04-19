package dto

import "certifisafe-back/model"

type PasswordRecoveryRequestDTO struct {
	Email string
}

type PasswordResetDTO struct {
	VerificationCode string
	NewPassword      string
}

func PasswordResetDTOtoModel(dto *PasswordResetDTO) *model.PasswordRecovery {
	return &model.PasswordRecovery{
		Code:        dto.VerificationCode,
		NewPassword: dto.NewPassword,
	}
}
