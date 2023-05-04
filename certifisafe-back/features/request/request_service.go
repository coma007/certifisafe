package request

import (
	certificate2 "certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"time"
)

type RequestService interface {
	CreateRequest(req *NewRequestDTO) (*RequestDTO, error)
	GetRequest(id int) (*RequestDTO, error)
	GetAllRequests() ([]*RequestDTO, error)
	GetAllRequestsByUser(userId int) ([]*RequestDTO, error)
	UpdateRequest(req *Request) error
	DeleteRequest(id int) error
	AcceptRequest(id int) error
	DeclineRequest(id int) error
}

type DefaultRequestService struct {
	repository         *DefaultRequestRepository
	certificateService *certificate2.DefaultCertificateService
}

func NewDefaultRequestService(repo *DefaultRequestRepository, certificateService *certificate2.DefaultCertificateService) *DefaultRequestService {
	return &DefaultRequestService{repo, certificateService}
}

func (service *DefaultRequestService) CreateRequest(req *NewRequestDTO) (*RequestDTO, error) {
	parentSerial := uint64(*req.ParentSerial)
	newRequest := Request{
		Datetime:            time.Time{},
		Status:              RequestStatus(PENDING),
		CertificateName:     req.CertificateName,
		CertificateType:     req.CertificateType,
		ParentCertificateID: &parentSerial,
		ParentCertificate:   certificate2.Certificate{},
		SubjectID:           req.SubjectId,
		Subject:             user.User{},
	}
	request, err := service.repository.CreateRequest(&newRequest)
	return RequestToDTO(request), err
}

func (service *DefaultRequestService) GetRequest(id int) (*RequestDTO, error) {
	request, err := service.repository.GetRequest(id)
	return RequestToDTO(request), err
}

func (service *DefaultRequestService) GetAllRequests() ([]*RequestDTO, error) {
	requests, err := service.repository.GetAllRequests()
	var requestsDTO []*RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *DefaultRequestService) GetAllRequestsByUser(userId int) ([]*RequestDTO, error) {
	requests, err := service.repository.GetAllRequestsByUser(userId)
	var requestsDTO []*RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *DefaultRequestService) UpdateRequest(req *Request) error {
	return service.repository.UpdateRequest(req)
}

func (service *DefaultRequestService) DeleteRequest(id int) error {
	return service.repository.DeleteRequest(id)
}

func (service *DefaultRequestService) AcceptRequest(id int) error {
	request, err := service.repository.GetRequest(id)
	if err != nil {
		return err
	}
	request.Status = ACCEPTED
	// TODO create certificate here
	return service.repository.UpdateRequest(request)
}

func (service *DefaultRequestService) DeclineRequest(id int) error {
	request, err := service.repository.GetRequest(id)
	if err != nil {
		return err
	}
	request.Status = REJECTED
	return service.repository.UpdateRequest(request)
}
