package repository

import (
	"bytes"
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/pavlo-v-chernykh/keystore-go/v4"
	"log"
	"math/big"
	"os"
	"time"
)

const store = "keystore.jsk"

var (
	ErrCertificateNotFound = errors.New("FromRepository - certificate not found")
)

type ICertificateRepository interface {

	UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error)
	GetCertificate(id int32) (model.Certificate, error)
	DeleteCertificate(id int32) error
	CreateCertificate(serialNumber big.Int, certPEM bytes.Buffer, certPrivKeyPEM bytes.Buffer) (x509.Certificate, error)
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

func (i *InmemoryCertificateRepository) CreateCertificate(serialNumber big.Int, certPEM bytes.Buffer,
	certPrivKeyPEM bytes.Buffer) (x509.Certificate, error) {
	config := utils.Config()
	password := []byte(config["keystore-password"])
	defer utils.Zeroing(password)

	pkeIn := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   certPrivKeyPEM.Bytes(),
		CertificateChain: []keystore.Certificate{
			{
				Type:    "X509",
				Content: certPEM.Bytes(),
			},
		},
	}
	ks := keystore.New()
	if err := ks.SetPrivateKeyEntry(fmt.Sprint(serialNumber), pkeIn, password); err != nil {
		return x509.Certificate{}, err
	}

	writeKeyStore(ks, store, password)

	ks = keystore.New()
	ks = readKeyStore(store, password)
	certificate, err := ks.GetPrivateKeyEntry(fmt.Sprint(serialNumber), password)
	if err != nil {
		return x509.Certificate{}, err
	}

	block, _ := pem.Decode(certificate.CertificateChain[0].Content)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return x509.Certificate{}, err
	}

	return *cert, nil
}

func readKeyStore(filename string, password []byte) keystore.KeyStore {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ks := keystore.New()
	if err := ks.Load(f, password); err != nil {
		log.Fatal(err) // nolint: gocritic
	}

	return ks
}

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = ks.Store(f, password)
	if err != nil {
		log.Fatal(err) // nolint: gocritic
	}
}
