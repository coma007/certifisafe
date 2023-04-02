package service

import (
	"certifisafe-back/model"
	"certifisafe-back/repository"
)

type RequestService interface {
	GetRequest(id int) (*model.Request, error)
	GetAllRequests() ([]*model.Request, error)
	CreateRequest(req *model.Request) (*model.Request, error)
	UpdateRequest(req *model.Request) error
	DeleteRequest(id int) error
}

type RequestServiceImpl struct {
	repository *repository.RequestRepositoryImpl
}

func NewRequestServiceImpl(repo *repository.RequestRepositoryImpl) *RequestServiceImpl {
	return &RequestServiceImpl{repo}
}

func (service *RequestServiceImpl) GetRequest(id int) (*model.Request, error) {
	return service.repository.GetRequest(id)
}

func (service *RequestServiceImpl) GetAllRequests() ([]*model.Request, error) {
	return service.repository.GetAllRequests()
}

func (service *RequestServiceImpl) CreateRequest(req *model.Request) (*model.Request, error) {
	return service.repository.CreateRequest(req)
}

func (service *RequestServiceImpl) UpdateRequest(req *model.Request) error {
	return service.repository.UpdateRequest(req)
}

func (service *RequestServiceImpl) DeleteRequest(id int) error {
	return service.repository.DeleteRequest(id)
}
