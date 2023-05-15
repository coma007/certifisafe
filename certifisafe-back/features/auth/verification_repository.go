package auth

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNoVerificationWithEmail = errors.New("no request for given email")
	ErrNoVerificationWithCode  = errors.New("no request for given code")
)

type VerificationRepository interface {
	GetVerification(id int32) (Verification, error)
	DeleteVerification(id int32) error
	CreateVerification(id int32, user Verification) (Verification, error)
	GetVerificationByCode(code string) (Verification, error)
	GetVerificationByEmail(email string) (*Verification, error)
	UpdateVerification(id int32, req Verification) (Verification, error)
}

type DefaultVerificationRepository struct {
	DB *gorm.DB
}

func NewDefaultVerificationRepository(db *gorm.DB) *DefaultVerificationRepository {
	return &DefaultVerificationRepository{
		DB: db,
	}
}

func (repository *DefaultVerificationRepository) CreateVerification(id int32, user Verification) (Verification, error) {
	result := repository.DB.Create(&user)
	return user, result.Error
}

func (repository *DefaultVerificationRepository) GetVerification(id int32) (Verification, error) {
	var verification Verification
	result := repository.DB.First(&verification, id)
	return verification, result.Error
}

func (repository *DefaultVerificationRepository) GetVerificationByCode(code string) (Verification, error) {
	var verification Verification
	result := repository.DB.Where("code=?", code).First(&verification)
	return verification, result.Error
}

func (repository *DefaultVerificationRepository) GetVerificationByEmail(email string) (*Verification, error) {
	var verification Verification
	result := repository.DB.Where("email=?", email).First(&verification)
	return &verification, result.Error
}

func (repository *DefaultVerificationRepository) UpdateVerification(id int32, req Verification) (Verification, error) {
	result := repository.DB.Save(&req)
	return req, result.Error
}

func (repository *DefaultVerificationRepository) DeleteVerification(id int32) error {
	result := repository.DB.Delete(&Verification{}, id)
	return result.Error
}
