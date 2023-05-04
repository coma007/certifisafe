package password_recovery

import "gorm.io/gorm"

type PasswordRecoveryRequest struct {
	gorm.Model
	Deleted gorm.DeletedAt
	Email   string
	Code    string
	IsUsed  bool
}

type PasswordRecovery struct {
	Code        string
	NewPassword string
}
