package model

type PasswordRecoveryRequest struct {
	Id    int
	Email string
	Code  string
}

type PasswordRecovery struct {
	Code        string
	NewPassword string
}
