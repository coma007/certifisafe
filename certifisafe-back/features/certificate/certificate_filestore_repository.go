package certificate

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

func (i *DefaultFileStoreCertificateRepository) CreateCertificate(serialNumber uint64, certPEM bytes.Buffer,
	certPrivKeyPEM bytes.Buffer) (x509.Certificate, error) {
	// TODO implement depending on chain

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
	certificate, err := i.GetCertificate(uint(serialNumber))
	if err != nil {
		return x509.Certificate{}, err
	}

	return certificate, nil
}

func (i *DefaultFileStoreCertificateRepository) GetCertificate(serialNumber uint) (x509.Certificate, error) {
	catls, err := tls.LoadX509KeyPair(getPublicName(uint64(serialNumber)), getPrivateName(uint64(serialNumber)))
	if err != nil {
		panic(err)
	}
	certificate, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		panic(err)
	}

	return *certificate, nil
}

func (i *DefaultFileStoreCertificateRepository) GetPrivateKey(serial uint) (*rsa.PrivateKey, error) {
	keyIn, err := os.ReadFile(getPrivateName(uint64(serial)))

	block, _ := pem.Decode(keyIn)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return &rsa.PrivateKey{}, err
	}

	return privateKey, nil
}

func getPrivateName(serial uint64) string {
	return "private" + string(os.PathSeparator) + strconv.FormatUint(serial, 10) + ".key"
}

func getPublicName(serial uint64) string {
	return "public" + string(os.PathSeparator) + strconv.FormatUint(serial, 10) + ".crt"
}
