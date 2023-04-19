package model

type PasswordRecoveryRequest struct {
	Id    int
	Email string
	Code  string
}
