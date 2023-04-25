package repository

import (
	"bytes"
	"certifisafe-back/model"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"math/big"
	"os"
)

type InmemoryKeyStoreCertificateRepository struct {
	Certificates []model.Certificate
	DB           *sql.DB
}

type IKeyStoreCertificateRepository interface {
	GetCertificate(id big.Int) (x509.Certificate, error)
	DeleteCertificate(id big.Int) error
	CreateCertificate(serialNumber big.Int, certPEM bytes.Buffer, certPrivKeyPEM bytes.Buffer) (x509.Certificate, error)
	GetPrivateKey(serial big.Int) (*rsa.PrivateKey, error)
}

func NewInMemoryCertificateKeyStoreRepository(db *sql.DB) *InmemoryKeyStoreCertificateRepository {
	var certificates = []model.Certificate{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
	}

	return &InmemoryKeyStoreCertificateRepository{
		Certificates: certificates,
		DB:           db,
	}
}

func (i *InmemoryKeyStoreCertificateRepository) GetCertificate(serialNumber big.Int) (x509.Certificate, error) {
	catls, err := tls.LoadX509KeyPair(getPublicName(serialNumber), getPrivateName(serialNumber))
	if err != nil {
		panic(err)
	}
	certificate, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}

	return *certificate, nil
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

func (i *InmemoryKeyStoreCertificateRepository) GetPrivateKey(serial big.Int) (*rsa.PrivateKey, error) {
	keyIn, err := os.ReadFile(getPrivateName(serial))

	block, _ := pem.Decode(keyIn)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	return privateKey, nil
}

func (i *InmemoryKeyStoreCertificateRepository) CreateCertificate(serialNumber big.Int, certPEM bytes.Buffer,
	certPrivKeyPEM bytes.Buffer) (x509.Certificate, error) {

	privateKey := certPrivKeyPEM.Bytes()
	publicKey := certPEM.Bytes()
	//certificateChain := []keystore.Certificate{
	//	{
	//		Type:    "X509",
	//		Content: certPEM.Bytes(),
	//	},
	//}
	certOut, err := os.Create(getPublicName(serialNumber))
	_, err = certOut.Write(publicKey)
	if err != nil {
		return x509.Certificate{}, err
	}
	defer certOut.Close()

	keyOut, err := os.OpenFile(getPrivateName(serialNumber), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	_, err = keyOut.Write(privateKey)
	if err != nil {
		return x509.Certificate{}, err
	}
	defer keyOut.Close()
	certificate, err := i.GetCertificate(serialNumber)
	if err != nil {
		return x509.Certificate{}, err
	}

	return certificate, nil

}

func getPrivateName(serial big.Int) string {
	return "private" + string(os.PathSeparator) + serial.String() + ".key"
}

func getPublicName(serial big.Int) string {
	return "public" + string(os.PathSeparator) + serial.String() + ".crt"
}
