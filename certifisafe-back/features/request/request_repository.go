package request

import (
	"certifisafe-back/features/certificate"
	user2 "certifisafe-back/features/user"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrRequestNotFound = errors.New("FromRepository - request not found")
)

type RequestRepository interface {
	CreateRequest(request *Request) (*Request, error)
	GetRequest(id int) (*Request, error)
	GetAllRequests() ([]*Request, error)
	GetAllRequestsByUser(userId int) ([]*Request, error)
	GetAllRequestsByUserSigning(userId int) ([]*Request, error)
	UpdateRequest(request *Request) error
	DeleteRequest(id int) error
}

type DefaultRequestRepository struct {
	DB                    *gorm.DB
	certificateRepository certificate.CertificateRepository
	userRepository        user2.UserRepository
}

func NewDefaultRequestRepository(db *gorm.DB, certificateRepo certificate.CertificateRepository, userRepo user2.UserRepository) *DefaultRequestRepository {
	return &DefaultRequestRepository{
		DB:                    db,
		certificateRepository: certificateRepo,
		userRepository:        userRepo,
	}
}

func (repository *DefaultRequestRepository) CreateRequest(request *Request) (*Request, error) {
	result := repository.DB.Create(&request)
	if result.Error != nil {
		return nil, result.Error
	}
	return repository.GetRequest(int(request.ID))
}

func (repository *DefaultRequestRepository) GetRequest(id int) (*Request, error) {
	request := &Request{}
	result := repository.DB.Preload("ParentCertificate").Preload("ParentCertificate.Issuer").
		Preload("ParentCertificate.Subject").Preload("Subject").Find(&request, id)
	return request, result.Error
}

func (repository *DefaultRequestRepository) GetAllRequests() ([]*Request, error) {
	requests := []*Request{}
	result := repository.DB.Preload("ParentCertificate").Preload("Subject").Find(&requests)
	return requests, result.Error
}

func (repository *DefaultRequestRepository) GetAllRequestsByUser(userId int) ([]*Request, error) {
	requests := []*Request{}
	result := repository.DB.Preload("ParentCertificate").Preload("Subject").Where("subject_id=?", userId).Find(&requests)
	return requests, result.Error
}

func (repository *DefaultRequestRepository) GetAllRequestsByUserSigning(userId int) ([]*Request, error) {
	requests := []*Request{}
	result := repository.DB.Preload("ParentCertificate").Where("parent_certificates.issuer_id = ?", userId).Find(&requests)
	return requests, result.Error
}

func (repository *DefaultRequestRepository) UpdateRequest(request *Request) error {
	result := repository.DB.Save(&request)
	return result.Error
}

func (repository *DefaultRequestRepository) DeleteRequest(id int) error {
	result := repository.DB.Delete(&Request{}, id)
	return result.Error
}
