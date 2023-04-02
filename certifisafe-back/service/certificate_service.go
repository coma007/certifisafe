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
	// "errors"
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
	UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error)
	GetCertificate(id int32) (model.Certificate, error)
	DeleteCertificate(id int32) error
	CreateCertificate(certificate x509.Certificate) (x509.Certificate, error)
}

type DefaultCertificateService struct {
	certificateRepo repository.ICertificateRepository
}

func NewDefaultCertificateService(cRepo repository.ICertificateRepository) *DefaultCertificateService {
	return &DefaultCertificateService{
		certificateRepo: cRepo,
	}
}

func (d *DefaultCertificateService) UpdateCertificate(id int32, certificate model.Certificate) (model.Certificate, error) {
	// if id <= 0 {
	// 	return ErrIDIsNotValid
	// }

	// if movie.Title == "" {
	// 	return ErrTitleIsNotEmpty
	// }

	// err := d.certificateRepo.UpdateCertificate(id, certificate)
	// if errors.Is(err, repository.ErrCertificateNotFound) {
	// 	return ErrCertificateNotFound
	// }

	return model.Certificate{}, nil
}
func (d *DefaultCertificateService) GetCertificate(id int32) (model.Certificate, error) {
	certificate, err := d.certificateRepo.GetCertificate(id)
	return certificate, err
}
func (d *DefaultCertificateService) DeleteCertificate(id int32) error {

	return nil
}
func (d *DefaultCertificateService) CreateCertificate(certificate x509.Certificate) (x509.Certificate, error) {
	// CA, root
	ca := &x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
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

	cert := &x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
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

	// create buffer and fill it with encoded value
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

	verifySignature(caPEM, certPEM)

	return *cert, nil
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
