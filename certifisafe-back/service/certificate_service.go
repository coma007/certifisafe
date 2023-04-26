package service

import (
	"bytes"
	"certifisafe-back/dto"
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
	CreateCertificate(cert dto.NewRequestDTO) (dto.CertificateDTO, error)
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
	certificate, err := d.certificateRepo.GetCertificate(id)
	return certificate, err
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

func (d *DefaultCertificateService) CreateCertificate(cert dto.NewRequestDTO) (dto.CertificateDTO, error) {
	// creating of leaf node

	var certificate x509.Certificate
	var certificatePEM bytes.Buffer
	var certificatePrivKeyPEM bytes.Buffer
	var err error

	subject := pkix.Name{
		CommonName:    cert.Certificate.Name,
		Organization:  []string{cert.Certificate.Name},
		Country:       []string{cert.Certificate.Name},
		StreetAddress: []string{cert.Certificate.Name},
		PostalCode:    []string{cert.Certificate.Name},
	}
	var parent x509.Certificate

	//chain
	//conf := tls.Config { }
	//conf.RootCAs = x509.NewCertPool()
	//for _, cert := range certChain.Certificate {
	//	x509Cert, err := x509.ParseCertificate(cert)
	//	if err != nil {
	//		panic(err)
	//	}
	//	conf.RootCAs.AddCert(x509Cert)
	//}

	switch dto.StringToType(cert.Certificate.Type) {
	case model.ROOT:
		certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateRootCa(subject)
		if err != nil {
			return dto.CertificateDTO{}, err
		}
		break
	case model.INTERMEDIATE:
		{
			parentSerial := new(big.Int)
			parentSerial.SetString(cert.ParentCertificate.Serial, 10)
			parent, err = d.certificateKeyStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return dto.CertificateDTO{}, err
			}

			privateKey, err := d.certificateKeyStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateSubordinateCa(subject, &parent, privateKey)
			if err != nil {
				return dto.CertificateDTO{}, err
			}
		}
	case model.END:
		{
			parentSerial := new(big.Int)
			parentSerial.SetString(cert.ParentCertificate.Serial, 10)
			parent, err = d.certificateKeyStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return dto.CertificateDTO{}, err
			}

			privateKey, err := d.certificateKeyStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateLeafCert(subject, &parent, privateKey)
			if err != nil {
				return dto.CertificateDTO{}, err
			}
		}
	default:
		{
			return dto.CertificateDTO{}, errors.New("invalid type of certificate given, try END, INTERMEDIATE or ROOT")
		}
	}

	certificateKeyStore, err := d.certificateKeyStoreRepo.CreateCertificate(*certificate.SerialNumber, certificatePEM, certificatePrivKeyPEM)
	if err != nil {
		return dto.CertificateDTO{}, err
	}

	certificateDB := model.Certificate{
		certificate.SerialNumber.String(),
		certificate.Subject.CommonName,
		// TODO fix nil values
		nil,
		nil,
		certificate.NotBefore,
		certificate.NotAfter,
		model.NOT_ACTIVE,
		dto.StringToType(cert.Certificate.Type),
	}

	certificateDB, err = d.certificateRepo.CreateCertificate(certificateDB)

	return *dto.X509CertificateToCertificateDTO(&certificateKeyStore), nil
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
