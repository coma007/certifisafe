package repository

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type InmemoryKeyStoreCertificateRepository struct {
	DB *gorm.DB
}

type IKeyStoreCertificateRepository interface {
	GetCertificate(id uint64) (x509.Certificate, error)
	DeleteCertificate(id uint64) error
	CreateCertificate(serialNumber uint64, certPEM bytes.Buffer, certPrivKeyPEM bytes.Buffer) (x509.Certificate, error)
	GetPrivateKey(serial uint64) (*rsa.PrivateKey, error)
}

func NewInMemoryCertificateKeyStoreRepository() *InmemoryKeyStoreCertificateRepository {
	return &InmemoryKeyStoreCertificateRepository{}
}

func (i *InmemoryKeyStoreCertificateRepository) GetCertificate(serialNumber uint64) (x509.Certificate, error) {
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

func (i *InmemoryKeyStoreCertificateRepository) DeleteCertificate(id uint64) error {
	//if i.Certificates[k].Id == id {
	//	// i.Certificates[k].Title = movie.Title
	//	return nil
	//}

	return nil
	//return ErrMovieNotFound
}

func (i *InmemoryKeyStoreCertificateRepository) GetPrivateKey(serial uint64) (*rsa.PrivateKey, error) {
	keyIn, err := os.ReadFile(getPrivateName(serial))

	block, _ := pem.Decode(keyIn)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	return privateKey, nil
}

func (i *InmemoryKeyStoreCertificateRepository) CreateCertificate(serialNumber uint64, certPEM bytes.Buffer,
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

func getPrivateName(serial uint64) string {
	return "private" + string(os.PathSeparator) + strconv.FormatUint(serial, 10) + ".key"
}

func getPublicName(serial uint64) string {
	return "public" + string(os.PathSeparator) + strconv.FormatUint(serial, 10) + ".crt"
}
