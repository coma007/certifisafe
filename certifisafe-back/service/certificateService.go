package service

import (
	"certifisafe-back/model"
	"certifisafe-back/repository"
	"errors"
	// "errors"
)

var (
	ErrIDIsNotValid               = errors.New("id is not valid")
	ErrCertificateNotFound        = errors.New("the certificate cannot be found")
	ErrIssuerNameIsNotValid       = errors.New("the certificate cannot be found")
	ErrFromIsNotValid             = errors.New("the certificate cannot be found")
	ErrToIsNotValid               = errors.New("the certificate cannot be found")
	ErrSubjectNameIsNotValid      = errors.New("the certificate cannot be found")
	ErrSubjectPublicKeyIsNotValid = errors.New("the certificate cannot be found")
	ErrIssuerIdIsNotValid         = errors.New("the certificate cannot be found")
	ErrSubjectIdIsNotValid        = errors.New("the certificate cannot be found")
	ErrSignatureIsNotValid        = errors.New("the certificate cannot be found")
)

type ICertificateService interface {
	UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error)
	GetCertificate(id int32) (model.Certificate, error)
	DeleteCertificate(id int32) error
	CreateCertificate(certificate model.Certificate) (model.Certificate, error)
}

type DefaultCertificateService struct {
	certificateRepo repository.ICertificateRepository
}

func NewDefaultCertificateService(cRepo repository.ICertificateRepository) *DefaultCertificateService {
	return &DefaultCertificateService{
		certificateRepo: cRepo,
	}
}

func (d *DefaultCertificateService) UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error) {
	// if id <= 0 {
	// 	return ErrIDIsNotValid
	// }

	// if movie.Title == "" {
	// 	return ErrTitleIsNotEmpty
	// }

	// err := d.certificateRepo.UpdateCertificate(id, certificate)
	// if errors.Is(err, repository.ErrCertificateNotFound) {
	// 	return ErrCertificateNotFound
	// }

	return model.Certificate{}, nil
}
func (d *DefaultCertificateService) GetCertificate(id int32) (model.Certificate, error) {

	return model.Certificate{}, nil
}
func (d *DefaultCertificateService) DeleteCertificate(id int32) error {

	return nil
}
func (d *DefaultCertificateService) CreateCertificate(certificate model.Certificate) (model.Certificate, error) {

	return model.Certificate{}, nil
}
