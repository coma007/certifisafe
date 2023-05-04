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

func NewInMemoryVerificationRepository(db *gorm.DB) *DefaultVerificationRepository {
	return &DefaultVerificationRepository{
		DB: db,
	}
}

func (i *DefaultVerificationRepository) CreateVerification(id int32, user Verification) (Verification, error) {
	result := i.DB.Create(user)
	return user, result.Error
}

func (i *DefaultVerificationRepository) GetVerification(id int32) (Verification, error) {
	var verification Verification
	result := i.DB.First(&verification, id)
	return verification, result.Error
}

func (i *DefaultVerificationRepository) GetVerificationByCode(code string) (Verification, error) {
	var verification Verification
	result := i.DB.Where("code=?", code).First(&verification)
	return verification, result.Error
}

func (i *DefaultVerificationRepository) GetVerificationByEmail(email string) (*Verification, error) {
	var verification Verification
	result := i.DB.Where("email=?", email).First(&verification)
	return &verification, result.Error
}

func (i *DefaultVerificationRepository) UpdateVerification(id int32, req Verification) (Verification, error) {
	result := i.DB.Save(req)
	return req, result.Error
}

func (i *DefaultVerificationRepository) DeleteVerification(id int32) error {
	result := i.DB.Delete(&Verification{}, id)
	return result.Error
}
