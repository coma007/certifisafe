package request

import (
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"errors"
	"time"
)

type RequestService interface {
	CreateRequest(req *NewRequestDTO) (*RequestDTO, error)
	GetRequest(id int) (*RequestDTO, error)
	GetAllRequests() ([]*RequestDTO, error)
	GetAllRequestsByUser(user user.User) ([]*RequestDTO, error)
	UpdateRequest(req *Request) error
	DeleteRequest(id int) error
	AcceptRequest(id int) error
	DeclineRequest(id int) error
}

type DefaultRequestService struct {
	repository         *DefaultRequestRepository
	certificateService *certificate.DefaultCertificateService
	userRepository     *user.DefaultUserRepository
}

func NewDefaultRequestService(repo *DefaultRequestRepository, certificateService *certificate.DefaultCertificateService, userRepository *user.DefaultUserRepository) *DefaultRequestService {
	return &DefaultRequestService{repo, certificateService, userRepository}
}

func (service *DefaultRequestService) CreateRequest(req *NewRequestDTO) (*RequestDTO, error) {
	subject, err := service.userRepository.GetUser(int32(req.SubjectId))
	if !subject.IsAdmin && req.CertificateType == certificate.TypeToString(certificate.ROOT) {
		return &RequestDTO{}, errors.New("cannot request for root certificate")
	}
	parentSerial := uint64(*req.ParentSerial)

	newRequest := Request{
		Datetime:            time.Time{},
		Status:              RequestStatus(PENDING),
		CertificateName:     req.CertificateName,
		CertificateType:     req.CertificateType,
		ParentCertificateID: &parentSerial,
		ParentCertificate:   certificate.Certificate{},
		SubjectID:           req.SubjectId,
		Subject:             user.User{},
	}
	request, err := service.repository.CreateRequest(&newRequest)
	service.acceptCertificateIfNeeded(request)
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

func (service *DefaultRequestService) GetAllRequestsByUser(user user.User) ([]*RequestDTO, error) {
	if user.IsAdmin {
		return service.GetAllRequests()
	}
	requests, err := service.repository.GetAllRequestsByUser(int(user.ID))
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

func (service *DefaultRequestService) acceptCertificateIfNeeded(request *Request) {
	parentCertificate, _ := service.certificateService.GetCertificate(*request.ParentCertificateID)
	if parentCertificate.Subject == request.Subject || request.Subject.IsAdmin {
		service.AcceptRequest(int(request.ID))
	}
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
