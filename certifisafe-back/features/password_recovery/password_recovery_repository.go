package password_recovery

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNoRequestWithEmail = errors.New("no request for given email")
	ErrNoRequestWithCode  = errors.New("no request for given code")
)

type PasswordRecoveryRepository interface {
	GetRequest(id int32) (PasswordRecoveryRequest, error)
	DeleteRequest(id int32) error
	UseRequestsForEmail(email string) error
	CreateRequest(id int32, user PasswordRecoveryRequest) (PasswordRecoveryRequest, error)
	GetRequestByCode(code string) (PasswordRecoveryRequest, error)
	GetRequestsByEmail(email string) ([]*PasswordRecoveryRequest, error)
	UpdateRequest(id int32, req PasswordRecoveryRequest) (PasswordRecoveryRequest, error)
}

type DefaultPasswordRecoveryRepository struct {
	DB *gorm.DB
}

func NewDefaultPasswordRecoveryRepository(db *gorm.DB) *DefaultPasswordRecoveryRepository {
	return &DefaultPasswordRecoveryRepository{
		DB: db,
	}
}

func (i *DefaultPasswordRecoveryRepository) CreateRequest(id int32, user PasswordRecoveryRequest) (PasswordRecoveryRequest, error) {
	result := i.DB.Create(&user)
	return user, result.Error
}

func (i *DefaultPasswordRecoveryRepository) GetRequest(id int32) (PasswordRecoveryRequest, error) {
	var r PasswordRecoveryRequest
	result := i.DB.First(&r, id)
	return r, result.Error
}

func (i *DefaultPasswordRecoveryRepository) GetRequestByCode(code string) (PasswordRecoveryRequest, error) {
	var r PasswordRecoveryRequest
	result := i.DB.Model(PasswordRecoveryRequest{}).Where("code=?", code).First(&r)

	return r, result.Error
}

func (i *DefaultPasswordRecoveryRepository) GetRequestsByEmail(email string) ([]*PasswordRecoveryRequest, error) {
	var requests []*PasswordRecoveryRequest
	result := i.DB.Where("email=?", email).Find(&requests)

	return requests, result.Error
}

func (i *DefaultPasswordRecoveryRepository) UpdateRequest(id int32, req PasswordRecoveryRequest) (PasswordRecoveryRequest, error) {
	result := i.DB.Save(&req)
	return req, result.Error
}

func (i *DefaultPasswordRecoveryRepository) DeleteRequest(id int32) error {
	result := i.DB.Delete(&PasswordRecoveryRequest{}, id)
	return result.Error
}

func (i *DefaultPasswordRecoveryRepository) UseRequestsForEmail(email string) error {
	requests, err := i.GetRequestsByEmail(email)
	if err != nil {
		return err
	}
	for j := 0; j < len(requests); j++ {
		if requests[j].IsUsed {
			continue
		}
		requests[j].IsUsed = true
		_, err = i.UpdateRequest(int32(requests[j].ID), *requests[j])
		if err != nil {
			return err
		}
	}
	return nil
}
