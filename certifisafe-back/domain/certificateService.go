package domain

import (
	"certifisafe-back/model"
	"certifisafe-back/repository"
	// "errors"
)

// var (
// 	ErrIDIsNotValid        = errors.New("id is not valid")
// 	ErrCertificateNotFound = errors.New("the certificate cannot be found")
// )

type ICertificateService interface {
	UpdateCertificate(id int, certificate model.Certificate) error
}

type DefaultCertificateService struct {
	certificateRepo repository.ICertificateRepository
}

func NewDefaultCertificateService(cRepo repository.ICertificateRepository) *DefaultCertificateService {
	return &DefaultCertificateService{
		certificateRepo: cRepo,
	}
}

func (d *DefaultCertificateService) UpdateCertificate(id int, certificate model.Certificate) error {
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

	return nil
}
