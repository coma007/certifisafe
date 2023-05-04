package certificate

import (
	"certifisafe-back/features/user"
	"gorm.io/gorm"
)

type CertificateRepository interface {
	CreateCertificate(certificate Certificate) (Certificate, error)
	GetCertificate(id uint64) (Certificate, error)
	GetCertificates() ([]Certificate, error)
	DeleteCertificate(id uint64) error
}

type DefaultCertificateRepository struct {
	DB *gorm.DB
}

func NewDefaultCertificateRepository(db *gorm.DB) *DefaultCertificateRepository {
	return &DefaultCertificateRepository{
		DB: db,
	}
}

func (i *DefaultCertificateRepository) CreateCertificate(certificate Certificate) (Certificate, error) {
	// TODO implement this

	//subject := 1
	//if certificate.Subject != nil {
	//	subject = certificate.Subject.Id
	//}
	//issuer := 1
	//
	//if certificate.Issuer != nil {
	//	issuer = certificate.Issuer.Id
	//}

	t := INTERMEDIATE
	empty := user.User{}
	if certificate.Issuer == empty {
		t = ROOT
	}
	//certificate.Subject = subject
	//certificate.Issuer = issuer
	certificate.Type = t

	result := i.DB.Create(&certificate)
	return certificate, result.Error
}

func (i *DefaultCertificateRepository) GetCertificate(id uint64) (Certificate, error) {
	var certificate Certificate
	result := i.DB.First(&certificate, id)
	return certificate, result.Error
}

func (i *DefaultCertificateRepository) GetCertificates() ([]Certificate, error) {
	var certificates []Certificate
	result := i.DB.Find(&certificates)
	return certificates, result.Error
}

func (i *DefaultCertificateRepository) DeleteCertificate(id uint64) error {
	result := i.DB.Delete(&Certificate{}, id)
	return result.Error
}
