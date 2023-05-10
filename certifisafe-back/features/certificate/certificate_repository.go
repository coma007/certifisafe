package certificate

import (
	"gorm.io/gorm"
)

type CertificateRepository interface {
	CreateCertificate(certificate Certificate) (Certificate, error)
	GetCertificate(id uint64) (Certificate, error)
	GetCertificates() ([]Certificate, error)
	GetLeafCertificates() ([]Certificate, error)
	UpdateCertificate(certificate *Certificate) error
	isRevoked(id uint64) (bool, error)
	BeginTransaction() *gorm.DB
	GetByUserId(id uint) ([]Certificate, error)
}

type DefaultCertificateRepository struct {
	DB *gorm.DB
}

func NewDefaultCertificateRepository(db *gorm.DB) *DefaultCertificateRepository {
	return &DefaultCertificateRepository{
		DB: db,
	}
}

func (i *DefaultCertificateRepository) CreateCertificate(certificate Certificate) (Certificate, error) {
	result := i.DB.Create(&certificate)
	if result.Error != nil {
		return Certificate{}, result.Error
	}
	return i.GetCertificate(uint64(certificate.ID))
}

func (i *DefaultCertificateRepository) GetCertificate(id uint64) (Certificate, error) {
	var certificate Certificate
	result := i.DB.Preload("Issuer").Preload("Subject").Preload("ParentCertificate").First(&certificate, id)
	return certificate, result.Error
}

func (i *DefaultCertificateRepository) GetCertificates() ([]Certificate, error) {
	var certificates []Certificate
	result := i.DB.Preload("Issuer").Preload("Subject").Find(&certificates)
	return certificates, result.Error
}

func (i *DefaultCertificateRepository) GetByUserId(id uint) ([]Certificate, error) {
	var certificates []Certificate
	result := i.DB.Preload("Issuer").Preload("Subject").Where("Subject.ID=?", id).Find(&certificates)
	return certificates, result.Error
}

func (i *DefaultCertificateRepository) GetLeafCertificates() ([]Certificate, error) {
	var certificates []Certificate
	result := i.DB.Where(
		"id NOT IN (?)",
		i.DB.Table("certificates").
			Select("parent_certificate_id").
			Where("parent_certificate_id IS NOT NULL")).
		Preload("ParentCertificate",
			func(db *gorm.DB) *gorm.DB {
				return db.Preload("ParentCertificate")
			}).
		Find(&certificates)

	return certificates, result.Error
}

func (i *DefaultCertificateRepository) UpdateCertificate(certificate *Certificate) error {
	result := i.DB.Save(&certificate)
	return result.Error
}

func (i *DefaultCertificateRepository) isRevoked(id uint64) (bool, error) {
	var count int64 = 1

	err := i.DB.
		Unscoped().
		Model(&Certificate{}).
		Where("status=?  and id=?", WITHDRAWN, id).
		Count(&count).
		Error
	return count != 0, err
}
func (i *DefaultCertificateRepository) BeginTransaction() *gorm.DB {
	return i.DB.Begin()
}
