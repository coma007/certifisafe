package repository

import (
	"certifisafe-back/model"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrRequestNotFound = errors.New("FromRepository - request not found")
)

type RequestRepository interface {
	GetRequest(id int) (*model.Request, error)
	GetAllRequests() ([]*model.Request, error)
	GetAllRequestsByUser() ([]*model.Request, error)
	CreateRequest(request *model.Request) (*model.Request, error)
	UpdateRequest(request *model.Request) error
	DeleteRequest(id int) error
}

type RequestRepositoryImpl struct {
	DB                    *gorm.DB
	certificateRepository ICertificateRepository
}

func NewRequestRepository(db *gorm.DB, certificateRepo ICertificateRepository) *RequestRepositoryImpl {
	return &RequestRepositoryImpl{
		DB:                    db,
		certificateRepository: certificateRepo,
	}
}

func (repository *RequestRepositoryImpl) GetRequest(id int) (*model.Request, error) {
	request := &model.Request{}
	result := repository.DB.Find(&request, id)
	return request, result.Error
}

func (repository *RequestRepositoryImpl) GetAllRequests() ([]*model.Request, error) {
	requests := []*model.Request{}
	result := repository.DB.Find(&requests)
	return requests, result.Error
}

func (repository *RequestRepositoryImpl) GetAllRequestsByUser(userId int) ([]*model.Request, error) {
	requests := []*model.Request{}
	result := repository.DB.Where("subject_id=?", userId).Find(&requests)
	return requests, result.Error
}

func (repository *RequestRepositoryImpl) CreateRequest(request *model.Request) (*model.Request, error) {
	result := repository.DB.Create(&request)
	return request, result.Error
}

func (repository *RequestRepositoryImpl) UpdateRequest(request *model.Request) error {
	result := repository.DB.Save(&request)
	return result.Error
}

func (repository *RequestRepositoryImpl) DeleteRequest(id int) error {
	result := repository.DB.Delete(&model.Request{}, id)
	return result.Error
}
