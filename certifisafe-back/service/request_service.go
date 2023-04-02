package service

import (
	"certifisafe-back/model"
	"certifisafe-back/repository"
)

type ReqestService interface {
	UpdateRequest(id int32, request model.Request) (model.Request, error)
	GetRequest(id int32) (model.Request, error)
	DeleteRequest(id int32) error
	CreateRequest(request model.Request) (model.Request, error)
}

type RequestServiceImpl struct {
	repository repository.IRequestRepository
}

func NewRequestServiceImpl(repo repository.IRequestRepository) *RequestServiceImpl {
	return &RequestServiceImpl{
		repository: repo,
	}
}

func (d *RequestServiceImpl) CreateRequest(request model.Request) (model.Request, error) {
	return model.Request{}, nil
}

func (d *RequestServiceImpl) UpdateReqest(request model.Request) (model.Request, error) {
	return model.Request{}, nil
}

func (d *RequestServiceImpl) GetRequest(id int32) (model.Request, error) {
	request, err := d.GetRequest(id)
	return request, err
}

func (d *RequestServiceImpl) DeleteRequest(id int32) error {
	return nil
}
