package service

import (
	"bytes"
	"certifisafe-back/model"
	"certifisafe-back/repository"
	"certifisafe-back/utils"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
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
	GetCertificates() ([]model.Certificate, error)
	DeleteCertificate(id big.Int) error
	CreateCertificate(subject pkix.Name, parentSerial big.Int, certificateType model.CertificateType) (x509.Certificate, error)
	IsValid(id big.Int) (bool, error)
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

func (d *DefaultCertificateService) GetCertificates() ([]model.Certificate, error) {
	certificates, err := d.certificateRepo.GetCertificates()
	if err != nil {
		return nil, err
	}
	return certificates, nil
}

func (d *DefaultCertificateService) DeleteCertificate(id big.Int) error {

	return nil
}

func (d *DefaultCertificateService) CreateCertificate(subject pkix.Name, parentSerial big.Int, kind model.CertificateType) (x509.Certificate, error) {
	var cert x509.Certificate
	var certPEM bytes.Buffer
	var certPrivKeyPEM bytes.Buffer
	var err error

	switch kind {
	case model.ROOT:
		{
			cert, certPEM, certPrivKeyPEM, err = GenerateRootCa(subject)
			if err != nil {
				return x509.Certificate{}, err
			}
			break
		}
	case model.INTERMEDIATE:
		{
			parent, err := d.certificateKeyStoreRepo.GetCertificate(parentSerial)
			if err != nil {
				return x509.Certificate{}, err
			}

			privateKey, err := d.certificateKeyStoreRepo.GetKey(parentSerial)
			cert, certPEM, certPrivKeyPEM, err = GenerateSubordinateCa(subject, &parent, privateKey)
			if err != nil {
				return x509.Certificate{}, err
			}
		}
	case model.END:
		{
			parent, err := d.certificateKeyStoreRepo.GetCertificate(parentSerial)
			if err != nil {
				return x509.Certificate{}, err
			}

			privateKey, err := d.certificateKeyStoreRepo.GetKey(parentSerial)
			cert, certPEM, certPrivKeyPEM, err = GenerateLeafCert(subject, &parent, privateKey)
			if err != nil {
				return x509.Certificate{}, err
			}
		}
	}

	//certResponse, err := d.certificateRepo.CreateCertificate(*certModel)
	//if err != nil {
	//	return x509.Certificate{}, err
	//}

	createCertificate, err := d.certificateKeyStoreRepo.CreateCertificate(*cert.SerialNumber, certPEM, certPrivKeyPEM)
	if err != nil {
		return x509.Certificate{}, err
	}

	return createCertificate, nil
}

func (d *DefaultCertificateService) IsValid(id big.Int) (bool, error) {
	certificate, err := d.certificateKeyStoreRepo.GetCertificate(id)
	if err != nil {
		return false, nil
	}

	if !d.checkChain(certificate) {
		return false, nil
	}

	if certificate.NotAfter.Before(time.Now()) || certificate.NotAfter.Before(certificate.NotBefore) {
		return false, nil
	}

	return true, nil
}

func (d *DefaultCertificateService) checkChain(certificate x509.Certificate) bool {
	serial, err := utils.StringToBigInt(certificate.Issuer.SerialNumber)
	if err != nil {
		return false
	}
	// if it is root check if it is self-signed
	if certificate.IsCA && serial.Cmp(certificate.SerialNumber) != 0 {
		err := certificate.CheckSignatureFrom(&certificate)
		if err != nil {
			return false
		} else {
			return true
		}
	}

	parentSerial, err := utils.StringToBigInt(certificate.Issuer.SerialNumber)
	if err != nil {
		return false
	}
	parent, err := d.certificateKeyStoreRepo.GetCertificate(parentSerial)
	if err != nil {
		return false
	}

	err = certificate.CheckSignatureFrom(&parent)
	if err != nil {
		return false
	}
	return d.checkChain(parent)
}
