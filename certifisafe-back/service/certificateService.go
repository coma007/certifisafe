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
	ErrIssuerNameIsNotValid       = errors.New("the issuer name is not valid")
	ErrFromIsNotValid             = errors.New("the from time is not valid")
	ErrToIsNotValid               = errors.New("the to time is not valid")
	ErrSubjectNameIsNotValid      = errors.New("the subject name is not valid")
	ErrSubjectPublicKeyIsNotValid = errors.New("the subject public key is not valid")
	ErrIssuerIdIsNotValid         = errors.New("the issuer id is not valid")
	ErrSubjectIdIsNotValid        = errors.New("the subject id is not valid")
	ErrSignatureIsNotValid        = errors.New("the signature is not valid")
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
