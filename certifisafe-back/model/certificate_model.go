package model

import (
	"time"
)

type Certificate struct {
	Id        int64
	Name      string
	Issuer    *User
	Subject   *User
	ValidFrom time.Time
	ValidTo   time.Time
	Status    CertificateStatus
	Type      CertificateType
}

type CertificateType int64
type CertificateStatus int64

const (
	ROOT CertificateType = iota
	INTERMEDIATE
	END
)
const (
	ACTIVE CertificateStatus = iota
	EXPIRED
	WITHDRAWN
	NOT_ACTIVE
)
