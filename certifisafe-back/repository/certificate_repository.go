package repository

import (
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"database/sql"
	"errors"
	"math/big"
)

const store = "keystore.jsk"

var (
	ErrCertificateNotFound = errors.New("FromRepository - certificate not found")
)

type ICertificateRepository interface {
	GetCertificate(id big.Int) (model.Certificate, error)
	DeleteCertificate(id big.Int) error
	CreateCertificate(certificate model.Certificate) error
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

func (i *InmemoryCertificateRepository) GetCertificate(id big.Int) (model.Certificate, error) {
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

func (i *InmemoryCertificateRepository) DeleteCertificate(id big.Int) error {
	for k := 0; k < len(i.Certificates); k++ {
		//if i.Certificates[k].Id == id {
		//	// i.Certificates[k].Title = movie.Title
		//	return nil
		//}
	}

	return nil
	//return ErrMovieNotFound
}

func (i *InmemoryCertificateRepository) CreateCertificate(certificate model.Certificate) error {
	subject := 1
	if certificate.Subject != nil {
		subject = certificate.Subject.Id
	}
	issuer := 1

	if certificate.Issuer != nil {
		issuer = certificate.Issuer.Id
	}
	t := model.INTERMEDIATE
	if certificate.Issuer == nil {
		t = model.ROOT
	}

	err := i.DB.QueryRow(
		"INSERT INTO certificates(id, name, valid_from, valid_to, subject_id, issuer_id, type, status) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", certificate.Id, certificate.Name, certificate.ValidFrom, certificate.ValidTo, subject, issuer, t, model.NOT_ACTIVE)
	if err != nil {
		return err.Err()
	}

	return nil
}
