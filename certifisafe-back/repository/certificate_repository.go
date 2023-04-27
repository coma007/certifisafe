package repository

import (
	"certifisafe-back/model"
	"errors"
	"gorm.io/gorm"
)

const store = "keystore.jsk"

var (
	ErrCertificateNotFound = errors.New("FromRepository - certificate not found")
)

type ICertificateRepository interface {
	GetCertificate(id int64) (model.Certificate, error)
	DeleteCertificate(id int64) error
	CreateCertificate(certificate model.Certificate) (model.Certificate, error)
	GetCertificates() ([]model.Certificate, error)
}

type InmemoryCertificateRepository struct {
	DB *gorm.DB
}

func NewInMemoryCertificateRepository(db *gorm.DB) *InmemoryCertificateRepository {
	return &InmemoryCertificateRepository{
		DB: db,
	}
}

func (i *InmemoryCertificateRepository) GetCertificate(id int64) (model.Certificate, error) {
	//TODO add subject and issuer
	var certificate model.Certificate
	result := i.DB.First(&certificate, id)
	return certificate, result.Error
}

func (i *InmemoryCertificateRepository) GetCertificates() ([]model.Certificate, error) {
	var certificates []model.Certificate
	//TODO add subject and issuer
	result := i.DB.Find(&certificates)

	return certificates, result.Error
}

func (i *InmemoryCertificateRepository) DeleteCertificate(id int64) error {
	result := i.DB.Delete(&model.Certificate{}, id)
	return result.Error
}

func (i *InmemoryCertificateRepository) CreateCertificate(certificate model.Certificate) (model.Certificate, error) {
	//subject := 1
	//if certificate.Subject != nil {
	//	subject = certificate.Subject.Id
	//}
	//issuer := 1
	//
	//if certificate.Issuer != nil {
	//	issuer = certificate.Issuer.Id
	//}

	t := model.INTERMEDIATE
	empty := model.User{}
	if certificate.Issuer == empty {
		t = model.ROOT
	}
	//certificate.Subject = subject
	//certificate.Issuer = issuer
	certificate.Type = t

	result := i.DB.Create(&certificate)
	return certificate, result.Error
}
