package repository

import (
	"certifisafe-back/model"
	"errors"
)

var (
	ErrCertificateNotFound = errors.New("FromRepository - certificate not found")
)

type ICertificateRepository interface {
	UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error)
	GetCertificate(id int32) (model.Certificate, error)
	DeleteCertificate(id int32) error
	CreateCertificate(id int32, certificate model.Certificate) (model.Certificate, error)
}

type InmemoryCertificateRepository struct {
	Certificates []model.Certificate
}

func NewInMemoryCertificateRepository() *InmemoryCertificateRepository {
	var certificates = []model.Certificate{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InmemoryCertificateRepository{
		Certificates: certificates,
	}
}

func (i *InmemoryCertificateRepository) UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error) {
	for k := 0; k < len(i.Certificates); k++ {
		if i.Certificates[k].Id == id {
			// i.Certificates[k].Title = movie.Title
			return model.Certificate{}, nil
		}
	}

	return model.Certificate{}, nil
	//return ErrMovieNotFound
}

func (i *InmemoryCertificateRepository) GetCertificate(id int32) (model.Certificate, error) {
	for k := 0; k < len(i.Certificates); k++ {
		if i.Certificates[k].Id == id {
			// i.Certificates[k].Title = movie.Title
			return model.Certificate{}, nil
		}
	}

	return model.Certificate{}, nil
	//return ErrMovieNotFound
}

func (i *InmemoryCertificateRepository) DeleteCertificate(id int32) error {
	for k := 0; k < len(i.Certificates); k++ {
		if i.Certificates[k].Id == id {
			// i.Certificates[k].Title = movie.Title
			return nil
		}
	}

	return nil
	//return ErrMovieNotFound
}

func (i *InmemoryCertificateRepository) CreateCertificate(id int32, certificate model.Certificate) (model.Certificate, error) {
	for k := 0; k < len(i.Certificates); k++ {
		if i.Certificates[k].Id == id {
			// i.Certificates[k].Title = movie.Title
			return model.Certificate{}, nil
		}
	}

	return model.Certificate{}, nil
	//return ErrMovieNotFound
}
