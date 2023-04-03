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
	CreateCertificate(certificate model.Certificate) (model.Certificate, error)
	GetCertificates() ([]model.Certificate, error)
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

func (i *InmemoryCertificateRepository) GetCertificates() ([]model.Certificate, error) {
	var result []model.Certificate
	rows, err := i.DB.Query("SELECT name, valid_from, valid_to  FROM certificates")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var certificate model.Certificate
		rows.Scan(&certificate.Subject, &certificate.ValidFrom, &certificate.ValidTo)
		result = append(result, certificate)
	}
	utils.CheckError(err)

	return result, nil
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

func (i *InmemoryCertificateRepository) CreateCertificate(certificate model.Certificate) (model.Certificate, error) {
	stmt, err := i.DB.Prepare(
		"INSERT INTO certificates(name, valid_from, valid_to, subject_id, subject_pk, issuer_id)" +
			"VALUES($1, $2, $3, $4, $5, $6);")
	utils.CheckError(err)

	err = stmt.QueryRow(certificate.Id, certificate.ValidFrom, certificate.ValidTo, 1,
		certificate.PublicKey, 1).Scan()

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
		}
		return model.Certificate{}, err

	}
	return certificate, nil
}
