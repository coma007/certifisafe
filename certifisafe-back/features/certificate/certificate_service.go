package certificate

import (
	"bytes"
	user2 "certifisafe-back/features/user"
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

type CertificateService interface {
	CreateCertificate(parentSerial *uint, certificateName string, certificateType string, subjectId uint) (CertificateDTO, error)
	GetCertificate(id uint64) (Certificate, error)
	GetCertificates() ([]Certificate, error)
	WithdrarwCertificate(id uint64) error
	IsValid(id uint64) (bool, error)
}

type DefaultCertificateService struct {
	certificateRepo         CertificateRepository
	certificateKeyStoreRepo FileStoreCertificateRepository
	userRepo                user2.UserRepository
}

func NewDefaultCertificateService(cRepo CertificateRepository, cKSRepo FileStoreCertificateRepository,
	uRepo user2.UserRepository) *DefaultCertificateService {
	return &DefaultCertificateService{
		certificateRepo:         cRepo,
		certificateKeyStoreRepo: cKSRepo,
		userRepo:                uRepo,
	}
}

func (d *DefaultCertificateService) CreateCertificate(parentSerial *uint, certificateName string, certificateType string, subjectId uint) (CertificateDTO, error) {
	// creating of leaf node

	var certificate x509.Certificate
	var certificatePEM bytes.Buffer
	var certificatePrivKeyPEM bytes.Buffer
	var err error

	// TODO change in dto
	subject := pkix.Name{
		CommonName: certificateName,
		//Organization:  []string{cert.Certificate.Name},
		//Country:       []string{cert.Certificate.Name},
		//StreetAddress: []string{cert.Certificate.Name},
		//PostalCode:    []string{cert.Certificate.Name},
	}
	var parent x509.Certificate

	// TODO add chain ?
	//conf := tls.Config { }
	//conf.RootCAs = x509.NewCertPool()
	//for _, cert := range certChain.Certificate {
	//	x509Cert, err := x509.ParseCertificate(cert)
	//	if err != nil {
	//		panic(err)
	//	}
	//	conf.RootCAs.AddCert(x509Cert)
	//}

	// TODO get from parent
	issuer, err := d.userRepo.CreateUser(user2.User{
		Email:     "issuer",
		Password:  "asd",
		FirstName: "qwe",
		LastName:  "qwe",
		Phone:     "ertert",
		IsAdmin:   false,
	})
	utils.CheckError(err)

	// TODO grab user from request, but wait until that TODO is resolved
	newSubject, err := d.userRepo.CreateUser(user2.User{
		Email:     "subject",
		Password:  "asd",
		FirstName: "asd",
		LastName:  "ads",
		Phone:     "adw",
		IsAdmin:   false,
	})
	utils.CheckError(err)

	certificateDB := Certificate{
		Name:      certificateName,
		Issuer:    issuer,
		Subject:   newSubject,
		ValidFrom: certificate.NotBefore,
		ValidTo:   certificate.NotAfter,
		Status:    NOT_ACTIVE,
		Type:      StringToType(certificateType),

		//IssuerID:  &issuer.Id,
		//SubjectID: &newSubject.Id,
	}

	certificateDB, err = d.certificateRepo.CreateCertificate(certificateDB)
	utils.CheckError(err)

	// TODO fix if chain added
	switch StringToType(certificateType) {
	case ROOT:
		certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateRootCa(subject, uint64(certificateDB.ID))
		if err != nil {
			return CertificateDTO{}, err
		}
		break
	case INTERMEDIATE:
		{
			parent, err = d.certificateKeyStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return CertificateDTO{}, err
			}

			privateKey, err := d.certificateKeyStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateSubordinateCa(subject, uint64(certificateDB.ID), &parent, privateKey)
			if err != nil {
				return CertificateDTO{}, err
			}
		}
	case END:
		{
			parent, err = d.certificateKeyStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return CertificateDTO{}, err
			}

			privateKey, err := d.certificateKeyStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateLeafCert(subject, uint64(certificateDB.ID), &parent, privateKey)
			if err != nil {
				return CertificateDTO{}, err
			}
		}
	default:
		{
			return CertificateDTO{}, errors.New("invalid type of certificate given, try END, INTERMEDIATE or ROOT")
		}
	}
	certificate.SerialNumber = new(big.Int).SetUint64(uint64(certificateDB.ID))

	_, err = d.certificateKeyStoreRepo.CreateCertificate(certificate.SerialNumber.Uint64(), certificatePEM, certificatePrivKeyPEM)
	if err != nil {
		return CertificateDTO{}, err
	}

	return *ModelToCertificateDTO(&certificateDB), nil
}

func (d *DefaultCertificateService) GetCertificate(id uint64) (Certificate, error) {
	certificate, err := d.certificateRepo.GetCertificate(id)
	return certificate, err
}

func (d *DefaultCertificateService) GetCertificates() ([]Certificate, error) {
	certificates, err := d.certificateRepo.GetCertificates()
	if err != nil {
		return nil, err
	}
	return certificates, nil
}

func (d *DefaultCertificateService) WithdrawCertificate(id uint64) error {

	certificate, err := d.GetCertificate(id)
	if err != nil {
		return err
	}

	transaction := d.certificateRepo.BeginTransaction()
	certificate, err = d.invalidateCertificate(&certificate)
	if err != nil {
		transaction.Rollback()
		return err
	}
	err = d.invalidateCertificatesSignedBy(certificate.ID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	return nil
}

func (d *DefaultCertificateService) IsValid(id uint64) (bool, error) {
	// TODO implement
	certificate, err := d.certificateKeyStoreRepo.GetCertificate(uint(id))
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
	parent, err := d.certificateKeyStoreRepo.GetCertificate(uint(parentSerial.Uint64()))
	if err != nil {
		return false
	}

	err = certificate.CheckSignatureFrom(&parent)
	if err != nil {
		return false
	}
	return d.checkChain(parent)
}

func (d *DefaultCertificateService) invalidateCertificate(certificate *Certificate) (Certificate, error) {
	certificate.ValidTo = time.Now()
	certificate.Status = CertificateStatus(WITHDRAWN)
	err := d.certificateRepo.UpdateCertificate(certificate)
	return *certificate, err
}

func (d *DefaultCertificateService) invalidateCertificatesSignedBy(serial uint) error {

	return nil
}
