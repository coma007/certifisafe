package certificate

import (
	"bytes"
	"certifisafe-back/features/user"
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
	CreateCertificate(parentSerial *uint, certificateName string, certificateType CertificateType, subjectId uint) (CertificateDTO, error)
	GetCertificate(id uint64) (Certificate, error)
	GetCertificates() ([]Certificate, error)
	IsValid(cert x509.Certificate) (bool, error)
	IsValidById(id uint64) (bool, error)
	GetCertificateFiles(certificateID uint64, user user.User) (string, string, error)
	WithdrawCertificate(certificateID uint64, user user.User) (CertificateDTO, error)
}

type DefaultCertificateService struct {
	certificateRepo          CertificateRepository
	certificateFileStoreRepo FileStoreCertificateRepository
	userRepo                 user.UserRepository
}

func NewDefaultCertificateService(certificateRepo CertificateRepository, fileStoreCertificateRepo FileStoreCertificateRepository,
	userRepository user.UserRepository) *DefaultCertificateService {
	return &DefaultCertificateService{
		certificateRepo:          certificateRepo,
		certificateFileStoreRepo: fileStoreCertificateRepo,
		userRepo:                 userRepository,
	}
}

func (service *DefaultCertificateService) CreateCertificate(parentSerial *uint, certificateName string, certificateType CertificateType, subjectId uint) (CertificateDTO, error) {
	var certificate x509.Certificate
	var certificatePEM bytes.Buffer
	var certificatePrivKeyPEM bytes.Buffer
	var err error

	var issuer user.User
	var parentCertificate *Certificate

	if parentSerial != nil {
		temp, err := service.certificateRepo.GetCertificate(uint64(*parentSerial))
		parentCertificate = &temp
		utils.CheckError(err)
		issuer = parentCertificate.Subject
	}

	newSubject, err := service.userRepo.GetUser(subjectId)
	utils.CheckError(err)

	certificateDB := Certificate{
		Name:              certificateName,
		Issuer:            issuer,
		Subject:           newSubject,
		Status:            ACTIVE,
		Type:              certificateType,
		ParentCertificate: parentCertificate,
	}

	subject := pkix.Name{
		CommonName: certificateName,
	}
	var parent x509.Certificate

	switch certificateType {
	case ROOT:
		{
			certificateDB, err = service.setDatesAndSave(&certificateDB, 5)
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
			certificateDB, err = service.setDatesAndSave(&certificateDB, 1)
			if err != nil {
				return CertificateDTO{}, err
			}

			parent, err = service.certificateFileStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return CertificateDTO{}, err
			}

			subject.SerialNumber = strconv.FormatUint(uint64(certificateDB.ID), 10)
			privateKey, err := service.certificateFileStoreRepo.GetPrivateKey(*parentSerial)
			certificate, certificatePEM, certificatePrivKeyPEM, err = GenerateSubordinateCa(subject, parent.Subject, uint64(certificateDB.ID), &parent, privateKey)
			if err != nil {
				return CertificateDTO{}, err
			}
		}
	case END:
		{
			certificateDB, err = service.setDatesAndSave(&certificateDB, 1)
			if err != nil {
				return CertificateDTO{}, err
			}

			parent, err = service.certificateFileStoreRepo.GetCertificate(*parentSerial)
			if err != nil {
				return CertificateDTO{}, err
			}

			subject.SerialNumber = strconv.FormatUint(uint64(certificateDB.ID), 10)
			privateKey, err := service.certificateFileStoreRepo.GetPrivateKey(*parentSerial)
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

	_, err = service.certificateFileStoreRepo.CreateCertificate(certificate.SerialNumber.Uint64(), certificatePEM, certificatePrivKeyPEM)
	if err != nil {
		return CertificateDTO{}, err
	}

	return *CertificateToDTO(&certificateDB), nil
}

func (service *DefaultCertificateService) setDatesAndSave(certificateDB *Certificate, years int) (Certificate, error) {
	validFrom := time.Now()
	validTo := time.Now().AddDate(years, 0, 0)
	certificateDB.ValidTo = validTo
	certificateDB.ValidFrom = validFrom

	return service.certificateRepo.CreateCertificate(*certificateDB)
}

func (service *DefaultCertificateService) GetCertificate(id uint64) (Certificate, error) {
	certificate, err := service.certificateRepo.GetCertificate(id)
	return certificate, err
}

func (service *DefaultCertificateService) GetCertificates() ([]Certificate, error) {
	certificates, err := service.certificateRepo.GetCertificates()
	if err != nil {
		return nil, err
	}
	return certificates, nil
}

func (service *DefaultCertificateService) GetCertificateFiles(certificateID uint64, user user.User) (string, string, error) {
	certificate, err := service.GetCertificate(certificateID)
	if err != nil {
		return "", "", err
	}
	var public, private string
	public = GetPublicName(certificateID)
	if user.IsAdmin || *certificate.SubjectID == int64(user.ID) {
		private = GetPrivateName(certificateID)
	}
	return public, private, err
}

func (service *DefaultCertificateService) WithdrawCertificate(id uint64, user user.User) (CertificateDTO, error) {

	certificate, err := service.GetCertificate(id)
	if err != nil {
		return CertificateDTO{}, err
	}
	//if !user.IsAdmin && *certificate.SubjectID != int64(user.ID) {
	//	// TODO also check if certificate has been withdrawn before
	//	return CertificateDTO{}, errors.New("no permissions")
	//}

	transaction := service.certificateRepo.BeginTransaction()
	{
		certificate, err = service.invalidateCertificate(&certificate)
		if err != nil {
			transaction.Rollback()
			return CertificateDTO{}, err
		}
		err = service.invalidateCertificatesSignedBy(&certificate)
		if err != nil {
			transaction.Rollback()
			return CertificateDTO{}, err
		}
	}
	transaction.Commit()
	return *CertificateToDTO(&certificate), nil
}

func (service *DefaultCertificateService) IsValidById(id uint64) (bool, error) {
	certificate, err := service.certificateFileStoreRepo.GetCertificate(uint(id))
	if err != nil {
		return false, err
	}
	return service.IsValid(certificate)
}

func (service *DefaultCertificateService) IsValid(certificate x509.Certificate) (bool, error) {

	if !isTimeValid(certificate) {
		return false, nil
	}

	isRevoked, err := service.certificateRepo.isRevoked(certificate.SerialNumber.Uint64())
	if err != nil {
		return false, err
	}
	if isRevoked {
		return false, nil
	}

	if !service.checkChain(certificate) {
		return false, nil
	}

	return true, nil
}

func isTimeValid(c x509.Certificate) bool {
	return !(c.NotAfter.Before(time.Now()) || c.NotAfter.Before(c.NotBefore))
}

func (service *DefaultCertificateService) checkChain(certificate x509.Certificate) bool {
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
	parent, err := service.certificateFileStoreRepo.GetCertificate(uint(parentSerial.Uint64()))
	if err != nil {
		return false
	}

	if !isTimeValid(parent) {
		return false
	}

	isRevoked, err := service.certificateRepo.isRevoked(parentSerial.Uint64())
	if isRevoked {
		return false
	}

	err = certificate.CheckSignatureFrom(&parent)
	if err != nil {
		return false
	}
	return service.checkChain(parent)
}

func (service *DefaultCertificateService) invalidateCertificate(certificate *Certificate) (Certificate, error) {
	certificate.ValidTo = time.Now()
	certificate.Status = WITHDRAWN
	err := service.certificateRepo.UpdateCertificate(certificate)
	return *certificate, err
}

func (service *DefaultCertificateService) invalidateCertificatesSignedBy(invalidCertificate *Certificate) error {
	leafCertificates, err := service.certificateRepo.GetLeafCertificates()
	if err != nil {
		return err
	}
	for _, leafCertificate := range leafCertificates {
		err = service.invalidateChain(&leafCertificate, invalidCertificate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *DefaultCertificateService) getChainPartToInvalidCertificate(chain []*Certificate, certificate *Certificate, invalidCertificate *Certificate) []*Certificate {
	if certificate.ID == invalidCertificate.ID {
		return chain
	}
	chain = append(chain, certificate)
	if certificate.Type != ROOT {
		return service.getChainPartToInvalidCertificate(chain, certificate.ParentCertificate, invalidCertificate)
	}
	return nil
}

func (service *DefaultCertificateService) invalidateChain(leafCertificate *Certificate, invalidCertificate *Certificate) error {
	var chain []*Certificate
	chain = service.getChainPartToInvalidCertificate(chain, leafCertificate, invalidCertificate)
	var err error
	if chain == nil {
		return nil
	}
	for _, certificate := range chain {
		_, err = service.invalidateCertificate(certificate)
		if err != nil {
			return err
		}
	}
	return nil
}
