package request

import (
	"certifisafe-back/features/certificate"
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
	GetAllRequestsByUser() ([]*Request, error)
	UpdateRequest(request *Request) error
	DeleteRequest(id int) error
}

type DefaultRequestRepository struct {
	DB                    *gorm.DB
	certificateRepository certificate.CertificateRepository
}

func NewDefaultRequestRepository(db *gorm.DB, certificateRepo certificate.CertificateRepository) *DefaultRequestRepository {
	return &DefaultRequestRepository{
		DB:                    db,
		certificateRepository: certificateRepo,
	}
}

func (repository *DefaultRequestRepository) CreateRequest(request *Request) (*Request, error) {
	result := repository.DB.Create(&request)
	return request, result.Error
}

func (repository *DefaultRequestRepository) GetRequest(id int) (*Request, error) {
	request := &Request{}
	result := repository.DB.Preload("ParentCertificate").Preload("Subject").Find(&request, id)
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

func (repository *DefaultRequestRepository) UpdateRequest(request *Request) error {
	result := repository.DB.Save(&request)
	return result.Error
}

func (repository *DefaultRequestRepository) DeleteRequest(id int) error {
	// TODO add logical deleting
	result := repository.DB.Delete(&Request{}, id)
	return result.Error
}
