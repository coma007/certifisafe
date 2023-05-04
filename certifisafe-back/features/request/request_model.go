package request

import (
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"gorm.io/gorm"
	"time"
)

type Request struct {
	gorm.Model
	//Id                  *int `gorm:"autoIncrement;PRIMARY_KEY"`
	Deleted             gorm.DeletedAt
	Datetime            time.Time
	Status              RequestStatus
	CertificateName     string
	CertificateType     string
	SubjectID           uint
	Subject             user.User `gorm:"foreignKey:SubjectID;"`
	ParentCertificateID *uint64
	ParentCertificate   certificate.Certificate `gorm:"foreignKey:ParentCertificateID;"`
}

type RequestStatus int64

const (
	PENDING RequestStatus = iota
	ACCEPTED
	REJECTED
)
