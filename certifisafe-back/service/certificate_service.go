package service

import (
	"bytes"
	"certifisafe-back/model"
	"certifisafe-back/repository"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"
)

var (
	ErrIDIsNotValid               = errors.New("id is not valid")
	ErrCertificateNotFound        = errors.New("the certificate cannot be found")
	ErrIssuerNameIsNotValid       = errors.New("the issuer name is not valid")
	ErrFromIsNotValid             = errors.New("the from time is not valid")
	ErrToIsNotValid               = errors.New("the to time is not valid")
	ErrSubjectNameIsNotValid      = errors.New("the subject name is not valid")
	ErrSubjectPublicKeyIsNotValid = errors.New("the subject public key is not valid")
	ErrIssuerIdIsNotValid         = errors.New("the issuer id is not valid")
	ErrSubjectIdIsNotValid        = errors.New("the subject id is not valid")
	ErrSignatureIsNotValid        = errors.New("the signature is not valid")
)

type ICertificateService interface {
	GetCertificate(id big.Int) (model.Certificate, error)
	DeleteCertificate(id big.Int) error
	CreateCertificate(certificate x509.Certificate) (x509.Certificate, error)
	IsValid(id big.Int) error
}

type DefaultCertificateService struct {
	certificateRepo         repository.ICertificateRepository
	certificateKeyStoreRepo repository.IKeyStoreCertificateRepository
}

func NewDefaultCertificateService(cRepo repository.ICertificateRepository, cKSRepo repository.IKeyStoreCertificateRepository) *DefaultCertificateService {
	return &DefaultCertificateService{
		certificateRepo:         cRepo,
		certificateKeyStoreRepo: cKSRepo,
	}
}

func (d *DefaultCertificateService) GetCertificate(id big.Int) (model.Certificate, error) {
	//certificate, err := d.certificateRepo.GetCertificate(id)
	return model.Certificate{}, nil
}
func (d *DefaultCertificateService) DeleteCertificate(id big.Int) error {

	return nil
}
func (d *DefaultCertificateService) CreateCertificate(certificate x509.Certificate) (x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	// CA, root
	ca := &x509.Certificate{
		Version:      3,
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
			SerialNumber:  serialNumber.String(),
		},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// generate private key for CA (private key contains public)
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return x509.Certificate{}, err
	}

	// create CA root certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return x509.Certificate{}, err
	}

	// create encoder
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// encode private key
	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	//check if it already exists
	serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)

	cert := &x509.Certificate{
		Version:      3,
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		Issuer:             ca.Subject,
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
		IPAddresses:        []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:          time.Now(),
		NotAfter:           time.Now().AddDate(10, 0, 0),
		IsCA:               false,
		SubjectKeyId:       []byte{1, 2, 3, 4, 6},
		//remove last one
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageAny},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	//generate private key
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return x509.Certificate{}, err
	}

	// create certificate and sign it with CA key
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return x509.Certificate{}, err
	}

	//create buffer and fill it with encoded value
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	//sign certificate with CA key
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	certResponse, err := d.certificateRepo.CreateCertificate(*cert.SerialNumber, *certPEM, *certPrivKeyPEM)
	if err != nil {
		return x509.Certificate{}, err
	}
	//verifySignature(caPEM, certPEM)

	return certResponse, nil
}

func verifySignature(rootPEM *bytes.Buffer, certPEM *bytes.Buffer) {
	// First, create the set of root certificates
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPEM.Bytes())
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode(certPEM.Bytes())
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		//DNSName: "mail.google.com",
		Roots: roots,
	}

	newCert, err := cert.Verify(opts)
	if err != nil {
		panic("failed to verify certificate: " + err.Error())
	}
	fmt.Println(newCert)
}

func (d *DefaultCertificateService) IsValid(id big.Int) (bool, error) {
	//certificate, err := d.certificateKeyStoreRepo.GetCertificate(id)
	//if err != nil {
	//	return false, nil
	//}
	//
	//parentSerial, err := utils.StringToBigInt(certificate.Issuer.SerialNumber)
	//if err != nil {
	//	return false, nil
	//}
	//parent, err := d.certificateKeyStoreRepo.GetCertificate(parentSerial)
	//if err != nil {
	//	return false, err
	//}

	//// create encoder
	//certPEM := new(bytes.Buffer)
	//pem.Encode(certPEM, &pem.Block{
	//	Type:  "CERTIFICATE",
	//	Bytes: caBytes,
	//})
	//
	//// encode private key
	//caPrivKeyPEM := new(bytes.Buffer)
	//pem.Encode(caPrivKeyPEM, &pem.Block{
	//	Type:  "RSA PRIVATE KEY",
	//	Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	//})
	//
	//verifySignature(parent.)
	return true, nil
}
