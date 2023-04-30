package model

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	gorm.Model
	Id int `gorm:"autoIncrement;PRIMARY_KEY"`

	Datetime time.Time
	Status   RequestStatus

	ParentCertificateID *uint64
	CertificateID       *uint64

	ParentCertificate Certificate `gorm:"foreignKey:ParentCertificateID;"`
	Certificate       Certificate `gorm:"foreignKey:CertificateID;"`
}

type RequestStatus int64

const (
	PENDING RequestStatus = iota
	ACCEPTED
	REJECTED
)
