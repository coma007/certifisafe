package model

import "time"

type Certificate struct {
	Id               int32
	Serial           string
	IssuerName       string // ???
	ValidFrom        time.Time
	ValidTo          time.Time
	SubjectName      string // ???
	SubjectPublicKey string
	IssuerId         string
	SubjectId        string
	Signature        string
}
