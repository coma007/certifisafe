package request

import (
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"gorm.io/gorm"
	"time"
)

type Request struct {
	gorm.Model
	Deleted             gorm.DeletedAt
	Datetime            time.Time
	Status              RequestStatus
	CertificateName     string
	CertificateType     certificate.CertificateType
	SubjectID           uint
	Subject             user.User `gorm:"foreignKey:SubjectID;"`
	ParentCertificateID *uint64
	ParentCertificate   certificate.Certificate `gorm:"foreignKey:ParentCertificateID;"`

	RejectedReason *string
}

type RequestStatus int64

const (
	PENDING RequestStatus = iota
	ACCEPTED
	REJECTED
)
