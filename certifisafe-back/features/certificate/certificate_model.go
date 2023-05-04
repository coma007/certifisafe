package certificate

import (
	"certifisafe-back/features/user"
	"gorm.io/gorm"
	"time"
)

// TODO add organization, postal code, ... to Certificate

type Certificate struct {
	gorm.Model
	Name      string
	Deleted   gorm.DeletedAt
	Issuer    user.User `gorm:"foreignKey:IssuerID;references:ID"`
	Subject   user.User `gorm:"foreignKey:SubjectID;references:ID"`
	ValidFrom time.Time
	ValidTo   time.Time
	Status    CertificateStatus
	Type      CertificateType
	IssuerID  *int64
	SubjectID *int64
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
