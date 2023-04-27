package repository

import (
	"bytes"
	"certifisafe-back/model"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type InmemoryKeyStoreCertificateRepository struct {
	Certificates []model.Certificate
	DB           *gorm.DB
}

type IKeyStoreCertificateRepository interface {
	GetCertificate(id int64) (x509.Certificate, error)
	DeleteCertificate(id int64) error
	CreateCertificate(serialNumber int64, certPEM bytes.Buffer, certPrivKeyPEM bytes.Buffer) (x509.Certificate, error)
	GetPrivateKey(serial int64) (*rsa.PrivateKey, error)
}

func NewInMemoryCertificateKeyStoreRepository() *InmemoryKeyStoreCertificateRepository {
	var certificates = []model.Certificate{}

	return &InmemoryKeyStoreCertificateRepository{
		Certificates: certificates,
	}
}

func (i *InmemoryKeyStoreCertificateRepository) GetCertificate(serialNumber int64) (x509.Certificate, error) {
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

func (i *InmemoryKeyStoreCertificateRepository) DeleteCertificate(id int64) error {
	for k := 0; k < len(i.Certificates); k++ {
		//if i.Certificates[k].Id == id {
		//	// i.Certificates[k].Title = movie.Title
		//	return nil
		//}
	}

	return nil
	//return ErrMovieNotFound
}

func (i *InmemoryKeyStoreCertificateRepository) GetPrivateKey(serial int64) (*rsa.PrivateKey, error) {
	keyIn, err := os.ReadFile(getPrivateName(serial))

	block, _ := pem.Decode(keyIn)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	return privateKey, nil
}

func (i *InmemoryKeyStoreCertificateRepository) CreateCertificate(serialNumber int64, certPEM bytes.Buffer,
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

func getPrivateName(serial int64) string {
	return "private" + string(os.PathSeparator) + strconv.FormatInt(serial, 10) + ".key"
}

func getPublicName(serial int64) string {
	return "public" + string(os.PathSeparator) + strconv.FormatInt(serial, 10) + ".crt"
}
