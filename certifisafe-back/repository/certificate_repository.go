package repository

import (
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"database/sql"
	"errors"
)

var (
	ErrCertificateNotFound = errors.New("FromRepository - certificate not found")
)

type ICertificateRepository interface {
	UpdateCertificate(id int64, certificate model.Certificate) (model.Certificate, error)
	GetCertificate(id int64) (model.Certificate, error)
	DeleteCertificate(id int64) error
	CreateCertificate(id int64, certificate model.Certificate) (model.Certificate, error)
}

type InmemoryCertificateRepository struct {
	Certificates []model.Certificate
	DB           *sql.DB
}

func NewInMemoryCertificateRepository(db *sql.DB) *InmemoryCertificateRepository {
	var certificates = []model.Certificate{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InmemoryCertificateRepository{
		Certificates: certificates,
		DB:           db,
	}
}

func (i *InmemoryCertificateRepository) UpdateCertificate(id int64, certificate model.Certificate) (model.Certificate, error) {
	for k := 0; k < len(i.Certificates); k++ {
		//if i.Certificates[k].Id == id {
		//	// i.Certificates[k].Title = movie.Title
		//	return model.Certificate{}, nil
		//}
	}

	return model.Certificate{}, nil
	//return ErrMovieNotFound
}

func (i *InmemoryCertificateRepository) GetCertificate(id int64) (model.Certificate, error) {
	stmt, err := i.DB.Prepare("SELECT id FROM certificates WHERE id=$1")
	utils.CheckError(err)

	var certificate model.Certificate
	err = stmt.QueryRow(id).Scan(&certificate.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
		}
		return model.Certificate{}, err

	}
	return certificate, nil
}

func (i *InmemoryCertificateRepository) DeleteCertificate(id int64) error {
	for k := 0; k < len(i.Certificates); k++ {
		//if i.Certificates[k].Id == id {
		//	// i.Certificates[k].Title = movie.Title
		//	return nil
		//}
	}

	return nil
	//return ErrMovieNotFound
}

func (i *InmemoryCertificateRepository) CreateCertificate(id int64, certificate model.Certificate) (model.Certificate, error) {
	for k := 0; k < len(i.Certificates); k++ {
		//if i.Certificates[k].Id == id {
		//	// i.Certificates[k].Title = movie.Title
		//	return model.Certificate{}, nil
		//}
	}

	return model.Certificate{}, nil
}
