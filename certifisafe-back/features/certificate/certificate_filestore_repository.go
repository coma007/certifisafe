package certificate

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type FileStoreCertificateRepository interface {
	GetCertificate(id uint) (x509.Certificate, error)
	CreateCertificate(serialNumber uint64, certPEM bytes.Buffer, certPrivKeyPEM bytes.Buffer) (x509.Certificate, error)
	GetPrivateKey(serial uint) (*rsa.PrivateKey, error)
}

type DefaultFileStoreCertificateRepository struct {
	DB *gorm.DB
}

func NewDefaultFileStoreCertificateRepository() *DefaultFileStoreCertificateRepository {
	return &DefaultFileStoreCertificateRepository{}
}

func (repository *DefaultFileStoreCertificateRepository) CreateCertificate(serialNumber uint64, certPEM bytes.Buffer,
	certPrivKeyPEM bytes.Buffer) (x509.Certificate, error) {

	privateKey := certPrivKeyPEM.Bytes()
	publicKey := certPEM.Bytes()
	certOut, err := os.Create(GetPublicName(serialNumber))
	_, err = certOut.Write(publicKey)
	if err != nil {
		return x509.Certificate{}, err
	}
	defer certOut.Close()

	keyOut, err := os.OpenFile(GetPrivateName(serialNumber), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	_, err = keyOut.Write(privateKey)
	if err != nil {
		return x509.Certificate{}, err
	}
	defer keyOut.Close()
	certificate, err := repository.GetCertificate(uint(serialNumber))
	if err != nil {
		return x509.Certificate{}, err
	}

	return certificate, nil
}

func (repository *DefaultFileStoreCertificateRepository) GetCertificate(serialNumber uint) (x509.Certificate, error) {
	catls, err := tls.LoadX509KeyPair(GetPublicName(uint64(serialNumber)), GetPrivateName(uint64(serialNumber)))
	if err != nil {
		return x509.Certificate{}, errors.New("no such certificate exists")
	}
	certificate, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		return x509.Certificate{}, errors.New("no such certificate exists")
	}

	return *certificate, nil
}

func (repository *DefaultFileStoreCertificateRepository) GetPrivateKey(serial uint) (*rsa.PrivateKey, error) {
	keyIn, err := os.ReadFile(GetPrivateName(uint64(serial)))

	block, _ := pem.Decode(keyIn)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	return privateKey, nil
}

func GetPrivateName(serial uint64) string {
	return "private" + string(os.PathSeparator) + strconv.FormatUint(serial, 10) + ".key"
}

func GetPublicName(serial uint64) string {
	return "public" + string(os.PathSeparator) + strconv.FormatUint(serial, 10) + ".crt"
}
