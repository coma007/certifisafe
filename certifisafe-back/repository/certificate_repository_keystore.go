package repository

import (
	"bytes"
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/pavlo-v-chernykh/keystore-go/v4"
	"log"
	"math/big"
	"os"
	"time"
)

type InmemoryKeyStoreCertificateRepository struct {
	Certificates []model.Certificate
	DB           *sql.DB
}

type IKeyStoreCertificateRepository interface {
	GetCertificate(id big.Int) (x509.Certificate, error)
	DeleteCertificate(id big.Int) error
	CreateCertificate(serialNumber big.Int, certPEM bytes.Buffer, certPrivKeyPEM bytes.Buffer) (x509.Certificate, error)
	GetKey(serial big.Int) rsa.PrivateKey
}

func NewInMemoryCertificateKeyStoreRepository(db *sql.DB) *InmemoryKeyStoreCertificateRepository {
	var certificates = []model.Certificate{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InmemoryKeyStoreCertificateRepository{
		Certificates: certificates,
		DB:           db,
	}
}

func (i *InmemoryKeyStoreCertificateRepository) GetCertificate(serialNumber big.Int) (x509.Certificate, error) {
	config := utils.Config()
	password := []byte(config["keystore-password"])
	defer utils.Zeroing(password)

	ks := keystore.New()
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

func (i *InmemoryKeyStoreCertificateRepository) DeleteCertificate(id big.Int) error {
	for k := 0; k < len(i.Certificates); k++ {
		//if i.Certificates[k].Id == id {
		//	// i.Certificates[k].Title = movie.Title
		//	return nil
		//}
	}

	return nil
	//return ErrMovieNotFound
}

func (i *InmemoryKeyStoreCertificateRepository) CreateCertificate(serialNumber big.Int, certPEM bytes.Buffer,
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

func (i *InmemoryKeyStoreCertificateRepository) GetKey(serial big.Int) (rsa.PrivateKey, error) {
	config := utils.Config()
	password := []byte(config["keystore-password"])
	ks := keystore.New()
	ks = readKeyStore(store, password)
	certificate, err := ks.GetPrivateKeyEntry(fmt.Sprint(serial), password)
	if err != nil {
		return rsa.PrivateKey{}, err
	}

	block, _ := pem.Decode(certificate.PrivateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return rsa.PrivateKey{}, err
	}

	return *privateKey, nil
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
