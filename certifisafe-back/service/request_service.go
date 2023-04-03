package service

import (
	"certifisafe-back/dto"
	"certifisafe-back/model"
	"certifisafe-back/repository"
)

type RequestService interface {
	GetRequest(id int) (*dto.RequestDTO, error)
	GetAllRequests() ([]*dto.RequestDTO, error)
	GetAllRequestsByUser(userId int) ([]*dto.RequestDTO, error)
	CreateRequest(req *dto.NewRequestDTO) (*dto.RequestDTO, error)
	UpdateRequest(req *model.Request) error
	DeleteRequest(id int) error
	AcceptRequest(id int) error
	DeclineRequest(id int) error
}

type RequestServiceImpl struct {
	repository         *repository.RequestRepositoryImpl
	certificateService *DefaultCertificateService
}

func NewRequestServiceImpl(repo *repository.RequestRepositoryImpl, certificateService *DefaultCertificateService) *RequestServiceImpl {
	return &RequestServiceImpl{repo, certificateService}
}

func (service *RequestServiceImpl) GetRequest(id int) (*dto.RequestDTO, error) {
	request, err := service.repository.GetRequest(id)
	return dto.RequestToDTO(request), err
}

func (service *RequestServiceImpl) GetAllRequests() ([]*dto.RequestDTO, error) {
	requests, err := service.repository.GetAllRequests()
	var requestsDTO []*dto.RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, dto.RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *RequestServiceImpl) GetAllRequestsByUser(userId int) ([]*dto.RequestDTO, error) {
	requests, err := service.repository.GetAllRequestsByUser(userId)
	var requestsDTO []*dto.RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, dto.RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *RequestServiceImpl) CreateRequest(req *dto.NewRequestDTO) (*dto.RequestDTO, error) {
	request := dto.NewRequestDTOtoModel(req)
	// TODO create certificate
	//request.Certificate = service.certificateService.CreateCertificate(request.Certificate)
	request, err := service.repository.CreateRequest(request)
	return dto.RequestToDTO(request), err
}

func (service *RequestServiceImpl) UpdateRequest(req *model.Request) error {
	return service.repository.UpdateRequest(req)
}

func (service *RequestServiceImpl) DeleteRequest(id int) error {
	return service.repository.DeleteRequest(id)
}

func (service *RequestServiceImpl) AcceptRequest(id int) error {
	request, err := service.repository.GetRequest(id)
	if err != nil {
		return err
	}
	request.Status = model.ACCEPTED
	return service.repository.UpdateRequest(request)
}

func (service *RequestServiceImpl) DeclineRequest(id int) error {
	request, err := service.repository.GetRequest(id)
	if err != nil {
		return err
	}
	request.Status = model.REJECTED
	return service.repository.UpdateRequest(request)
}
