package repository

import (
	"certifisafe-back/model"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNoRequestWithEmail = errors.New("no request for given email")
	ErrNoRequestWithCode  = errors.New("no request for given code")
)

type IPasswordRecoveryRepository interface {
	GetRequest(id int32) (model.PasswordRecoveryRequest, error)
	DeleteRequest(id int32) error
	UseRequestsForEmail(email string) error
	CreateRequest(id int32, user model.PasswordRecoveryRequest) (model.PasswordRecoveryRequest, error)
	GetRequestByCode(code string) (model.PasswordRecoveryRequest, error)
	GetRequestsByEmail(email string) ([]*model.PasswordRecoveryRequest, error)
	UpdateRequest(id int32, req model.PasswordRecoveryRequest) (model.PasswordRecoveryRequest, error)
}

type InMemoryPasswordRecoveryRepository struct {
	Requests []model.PasswordRecoveryRequest
	DB       *gorm.DB
}

func NewInMemoryPasswordRecoveryRepository(db *gorm.DB) *InMemoryPasswordRecoveryRepository {
	var requests = []model.PasswordRecoveryRequest{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InMemoryPasswordRecoveryRepository{
		Requests: requests,
		DB:       db,
	}
}

func (i *InMemoryPasswordRecoveryRepository) GetRequest(id int32) (model.PasswordRecoveryRequest, error) {
	var r model.PasswordRecoveryRequest
	result := i.DB.First(&r, id)
	return r, result.Error
}

func (i *InMemoryPasswordRecoveryRepository) UpdateRequest(id int32, req model.PasswordRecoveryRequest) (model.PasswordRecoveryRequest, error) {
	result := i.DB.Save(req)
	return req, result.Error
}

func (i *InMemoryPasswordRecoveryRepository) DeleteRequest(id int32) error {
	result := i.DB.Delete(&model.Request{}, id)
	return result.Error
}

func (i *InMemoryPasswordRecoveryRepository) CreateRequest(id int32, user model.PasswordRecoveryRequest) (model.PasswordRecoveryRequest, error) {
	result := i.DB.Create(user)
	return user, result.Error
}

func (i *InMemoryPasswordRecoveryRepository) GetRequestByCode(code string) (model.PasswordRecoveryRequest, error) {
	var r model.PasswordRecoveryRequest
	result := i.DB.Where("code=?", code).First(&r)

	return r, result.Error
}

func (i *InMemoryPasswordRecoveryRepository) GetRequestsByEmail(email string) ([]*model.PasswordRecoveryRequest, error) {
	var requests []*model.PasswordRecoveryRequest
	result := i.DB.Where("email=?", email).Find(requests)

	return requests, result.Error
}

func (i *InMemoryPasswordRecoveryRepository) UseRequestsForEmail(email string) error {
	requests, err := i.GetRequestsByEmail(email)
	if err != nil {
		return err
	}
	for j := 0; j < len(requests); j++ {
		if requests[j].IsUsed {
			continue
		}
		requests[j].IsUsed = true
		_, err = i.UpdateRequest(int32(requests[j].Id), *requests[j])
		if err != nil {
			return err
		}
	}
	return nil
}
