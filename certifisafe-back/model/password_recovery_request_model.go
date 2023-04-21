package model

type PasswordRecoveryRequest struct {
	Id     int
	Email  string
	Code   string
	IsUsed bool
}

type PasswordRecovery struct {
	Code        string
	NewPassword string
}
