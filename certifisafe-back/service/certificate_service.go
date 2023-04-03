package service

import (
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
	DeleteCertificate(id big.Int) error
	CreateCertificate(certificate x509.Certificate) error
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
func (d *DefaultCertificateService) DeleteCertificate(id big.Int) error {

	return nil
}

func (d *DefaultCertificateService) CreateCertificate(certificate x509.Certificate, parentSerial big.Int) (x509.Certificate, error) {
	// creating of leaf node
	parent, err := d.certificateKeyStoreRepo.GetCertificate(parentSerial)
	if err != nil {
		return err
	}
	subject := pkix.Name{
		Country:            nil,
		Organization:       nil,
		OrganizationalUnit: nil,
		PostalCode:         nil,
		CommonName:         "",
		Names:              nil,
	}
	cert, certPEM, certPrivKeyPEM, err := GenerateLeafCert(subject, &parent, d.certificateKeyStoreRepo.GetKey(parentSerial))
	if err != nil {
		return err
	}

	//certResponse, err := d.certificateRepo.CreateCertificate(*certModel)
	//if err != nil {
	//	return x509.Certificate{}, err
	//}

	createCertificate, err := d.certificateKeyStoreRepo.CreateCertificate(*cert.SerialNumber, certPEM, certPrivKeyPEM)
	if err != nil {
		return err
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

// TODO TEST
func (d *DefaultCertificateService) checkChain(certificate x509.Certificate) bool {
	if certificate.IsCA {
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
