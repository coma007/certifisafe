package certificate

import (
	"bytes"
	user2 "certifisafe-back/features/user"
	"certifisafe-back/utils"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"strconv"
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
	GetCertificate(id uint64) (Certificate, error)
	GetCertificates() ([]Certificate, error)
	DeleteCertificate(id uint64) error
	CreateCertificate(parentSerial *uint, certificateName string, certificateType string, subjectId uint) (CertificateDTO, error)
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
	var certificate x509.Certificate
	var certificatePEM bytes.Buffer
	var certificatePrivKeyPEM bytes.Buffer
	var err error

	var issuer user2.User
	var parentCertificate *Certificate

	if parentSerial != nil {
		temp, err := d.certificateRepo.GetCertificate(uint64(*parentSerial))
		parentCertificate = &temp
		utils.CheckError(err)
		issuer = parentCertificate.Subject
	}

	newSubject, err := d.userRepo.GetUser(subjectId)
	utils.CheckError(err)

	certificateDB := Certificate{
		Name:              certificateName,
		Issuer:            issuer,
		Subject:           newSubject,
		Status:            NOT_ACTIVE,
		Type:              StringToType(certificateType),
		ParentCertificate: parentCertificate,
	}

	subject := pkix.Name{
		CommonName: certificateName,
	}
	var parent x509.Certificate

	switch StringToType(certificateType) {
	case ROOT:
		{
			certificateDB, err = d.setDatesAndSave(&certificateDB, 5)
			if err != nil {
				return CertificateDTO{}, err
			}
			subject.SerialNumber = strconv.FormatUint(uint64(certificateDB.ID), 10)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateRootCa(subject, uint64(certificateDB.ID))

			if err != nil {
				return CertificateDTO{}, err
			}
			break
		}
	case INTERMEDIATE:
		{
			certificateDB, err = d.setDatesAndSave(&certificateDB, 1)
			if err != nil {
				return CertificateDTO{}, err
			}

			parent, err = d.certificateKeyStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return CertificateDTO{}, err
			}

			subject.SerialNumber = strconv.FormatUint(uint64(certificateDB.ID), 10)
			privateKey, err := d.certificateKeyStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateSubordinateCa(subject, parent.Subject, uint64(certificateDB.ID), &parent, privateKey)
			if err != nil {
				return CertificateDTO{}, err
			}
		}
	case END:
		{
			certificateDB, err = d.setDatesAndSave(&certificateDB, 1)
			if err != nil {
				return CertificateDTO{}, err
			}

			parent, err = d.certificateKeyStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return CertificateDTO{}, err
			}

			subject.SerialNumber = strconv.FormatUint(uint64(certificateDB.ID), 10)
			privateKey, err := d.certificateKeyStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateLeafCert(subject, parent.Subject, uint64(certificateDB.ID), &parent, privateKey)
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

func (d *DefaultCertificateService) setDatesAndSave(certificateDB *Certificate, years int) (Certificate, error) {
	validFrom := time.Now()
	validTo := time.Now().AddDate(years, 0, 0)
	certificateDB.ValidTo = validTo
	certificateDB.ValidFrom = validFrom

	return d.certificateRepo.CreateCertificate(*certificateDB)
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

func (d *DefaultCertificateService) DeleteCertificate(id uint64) error {
	return d.certificateRepo.DeleteCertificate(id)
}

func (d *DefaultCertificateService) IsValid(id uint64) (bool, error) {
	certificate, err := d.certificateKeyStoreRepo.GetCertificate(uint(id))
	if err != nil {
		return false, nil
	}

	if !isTimeValid(certificate) {
		return false, nil
	}

	isRevoked, err := d.certificateRepo.isRevoked(id)
	if isRevoked {
		return false, nil
	}

	if !d.checkChain(certificate) {
		return false, nil
	}

	return true, nil
}

func isTimeValid(c x509.Certificate) bool {
	return !(c.NotAfter.Before(time.Now()) || c.NotAfter.Before(c.NotBefore))
}

func (d *DefaultCertificateService) checkChain(certificate x509.Certificate) bool {
	// if it is root check if it is self-signed(subject is same as issuer)
	if certificate.IsCA && certificate.Issuer.SerialNumber == certificate.SerialNumber.String() {
		err := certificate.CheckSignatureFrom(&certificate)
		if err != nil {
			return false
		} else {
			return true
		}
	}

	parentSerial, err := utils.StringToBigInt(certificate.Issuer.SerialNumber)
	//parentSerial, err := utils.StringToBigInt("4")
	if err != nil {
		return false
	}
	parent, err := d.certificateKeyStoreRepo.GetCertificate(uint(parentSerial.Uint64()))
	if err != nil {
		return false
	}

	if !isTimeValid(parent) {
		return false
	}

	isRevoked, err := d.certificateRepo.isRevoked(parentSerial.Uint64())
	if isRevoked {
		return false
	}

	err = certificate.CheckSignatureFrom(&parent)
	if err != nil {
		return false
	}
	return d.checkChain(parent)
}
