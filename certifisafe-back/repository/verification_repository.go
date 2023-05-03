package repository

import (
	"certifisafe-back/model"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNoVerificationWithEmail = errors.New("no request for given email")
	ErrNoVerificationWithCode  = errors.New("no request for given code")
)

type IVerificationRepository interface {
	GetVerification(id int32) (model.Verification, error)
	DeleteVerification(id int32) error
	CreateVerification(id int32, user model.Verification) (model.Verification, error)
	GetVerificationByCode(code string) (model.Verification, error)
	GetVerificationByEmail(email string) (*model.Verification, error)
	UpdateVerification(id int32, req model.Verification) (model.Verification, error)
}

type InMemoryVerificationRepository struct {
	DB *gorm.DB
}

func NewInMemoryVerificationRepository(db *gorm.DB) *InMemoryVerificationRepository {

	return &InMemoryVerificationRepository{
		DB: db,
	}
}

func (i *InMemoryVerificationRepository) GetVerification(id int32) (model.Verification, error) {
	var verification model.Verification
	result := i.DB.First(&verification, id)

	return verification, result.Error
}

func (i *InMemoryVerificationRepository) UpdateVerification(id int32, req model.Verification) (model.Verification, error) {
	result := i.DB.Save(req)
	return req, result.Error
}

func (i *InMemoryVerificationRepository) DeleteVerification(id int32) error {
	result := i.DB.Delete(&model.Verification{}, id)
	return result.Error
}

func (i *InMemoryVerificationRepository) CreateVerification(id int32, user model.Verification) (model.Verification, error) {
	result := i.DB.Create(user)
	return user, result.Error
}

func (i *InMemoryVerificationRepository) GetVerificationByCode(code string) (model.Verification, error) {
	var verification model.Verification
	result := i.DB.Where("code=?", code).First(&verification)
	return verification, result.Error
}

func (i *InMemoryVerificationRepository) GetVerificationByEmail(email string) (*model.Verification, error) {
	var verification model.Verification
	result := i.DB.Where("email=?", email).First(&verification)
	return &verification, result.Error
}
