package request

import (
	"certifisafe-back/features/auth"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"errors"
	"time"
)

type RequestService interface {
	CreateRequest(req *NewRequestDTO, subject *user.User) (*RequestDTO, error)
	GetRequest(id int) (*RequestDTO, error)
	GetAllRequests() ([]*RequestDTO, error)
	GetAllRequestsByUserSigning(user user.User) ([]*RequestDTO, error)
	GetAllRequestsByUser(user user.User) ([]*RequestDTO, error)
	UpdateRequest(req *Request) error
	DeleteRequest(id int) error
	AcceptRequest(id int) (*Request, error)
	DeclineRequest(id int, reason string) error
}

type DefaultRequestService struct {
	requestRepository  RequestRepository
	certificateService certificate.CertificateService
	userRepo           user.UserRepository
	authService        auth.AuthService
}

func NewDefaultRequestService(requestRepo RequestRepository, certificateService certificate.CertificateService, userRepo user.UserRepository, authService auth.AuthService) *DefaultRequestService {
	return &DefaultRequestService{requestRepo, certificateService, userRepo, authService}
}

func (service *DefaultRequestService) CreateRequest(req *NewRequestDTO, subject *user.User) (*RequestDTO, error) {
	if req.ParentSerial == nil {
		req.CertificateType = "ROOT"
	}
	if !subject.IsAdmin && req.CertificateType == certificate.TypeToString(certificate.ROOT) {
		return &RequestDTO{}, errors.New("cannot request for root certificate")
	}
	parentSerial := uint64(*req.ParentSerial)

	newRequest := Request{
		Datetime:            time.Now(),
		Status:              RequestStatus(PENDING),
		CertificateName:     req.CertificateName,
		CertificateType:     certificate.StringToType(req.CertificateType),
		ParentCertificateID: &parentSerial,
		ParentCertificate:   certificate.Certificate{},
		SubjectID:           subject.ID,
		Subject:             user.User{},
	}
	request, err := service.requestRepository.CreateRequest(&newRequest)
	request, err = service.acceptCertificateIfNeeded(request)
	return RequestToDTO(request), err
}

func (service *DefaultRequestService) GetRequest(id int) (*RequestDTO, error) {
	request, err := service.requestRepository.GetRequest(id)
	return RequestToDTO(request), err
}

func (service *DefaultRequestService) GetAllRequests() ([]*RequestDTO, error) {
	requests, err := service.requestRepository.GetAllRequests()
	var requestsDTO []*RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *DefaultRequestService) GetAllRequestsByUserSigning(user user.User) ([]*RequestDTO, error) {
	if user.IsAdmin {
		return service.GetAllRequests()
	}
	requests, err := service.requestRepository.GetAllRequestsByUser(int(user.ID))
	var requestsDTO []*RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *DefaultRequestService) GetAllRequestsByUser(user user.User) ([]*RequestDTO, error) {
	requests, err := service.requestRepository.GetAllRequestsByUser(int(user.ID))
	var requestsDTO []*RequestDTO
	for i := 0; i < len(requests); i++ {
		requestsDTO = append(requestsDTO, RequestToDTO(requests[i]))
	}
	return requestsDTO, err
}

func (service *DefaultRequestService) UpdateRequest(req *Request) error {
	return service.requestRepository.UpdateRequest(req)
}

func (service *DefaultRequestService) DeleteRequest(id int) error {
	return service.requestRepository.DeleteRequest(id)
}

func (service *DefaultRequestService) acceptCertificateIfNeeded(request *Request) (*Request, error) {
	parentCertificate, _ := service.certificateService.GetCertificate(*request.ParentCertificateID)
	var err error = nil
	if parentCertificate.Subject.ID == request.Subject.ID || request.Subject.IsAdmin {
		request, err = service.AcceptRequest(int(request.ID))
	}
	return request, err
}

func (service *DefaultRequestService) AcceptRequest(id int) (*Request, error) {
	request, err := service.requestRepository.GetRequest(id)
	if err != nil {
		return nil, err
	}
	request.Status = ACCEPTED
	parentSerial := uint(*request.ParentCertificateID)
	_, err = service.certificateService.CreateCertificate(&parentSerial, request.CertificateName, request.CertificateType, request.SubjectID)
	if err != nil {
		return nil, err
	}
	return request, service.requestRepository.UpdateRequest(request)
}

func (service *DefaultRequestService) DeclineRequest(id int, reason string) error {
	request, err := service.requestRepository.GetRequest(id)
	if err != nil {
		return err
	}
	request.Status = REJECTED
	request.RejectedReason = &reason
	return service.requestRepository.UpdateRequest(request)
}
