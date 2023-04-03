package repository

import (
	"certifisafe-back/model"
	"database/sql"
	"errors"
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
	DB                    *sql.DB
	certificateRepository ICertificateRepository
}

func NewRequestRepository(db *sql.DB, certificateRepo ICertificateRepository) *RequestRepositoryImpl {
	return &RequestRepositoryImpl{
		DB:                    db,
		certificateRepository: certificateRepo,
	}
}

func (repository *RequestRepositoryImpl) GetRequest(id int) (*model.Request, error) {
	request := &model.Request{}

	var parentCertificateId int64
	var certificateId int64
	err := repository.DB.QueryRow("SELECT id, datetime, status, parent_certificate_id, certificate_id FROM requests WHERE id = $1", id).Scan(&request.Id, &request.Datetime, &request.Status, parentCertificateId, certificateId)
	//repository.certificateRepository.GetCertificate(parentCertificateId)
	//repository.certificateRepository.GetCertificate(certificateId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}
	return request, nil
}

func (repository *RequestRepositoryImpl) GetAllRequests() ([]*model.Request, error) {
	rows, err := repository.DB.Query("SELECT id, datetime, status FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []*model.Request{}

	for rows.Next() {
		r := &model.Request{}
		err := rows.Scan(&r.Id, &r.Datetime, &r.Status)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}

func (repository *RequestRepositoryImpl) GetAllRequestsByUser(userId int) ([]*model.Request, error) {
	rows, err := repository.DB.Query("SELECT id, datetime, status FROM requests WHERE subject_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []*model.Request{}

	for rows.Next() {
		r := &model.Request{}
		err := rows.Scan(&r.Id, &r.Datetime, &r.Status)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}

func (repository *RequestRepositoryImpl) CreateRequest(request *model.Request) (*model.Request, error) {
	var id int
	err := repository.DB.QueryRow("INSERT INTO requests(datetime, parent_certificate_pk, certificate_pk, status) VALUES($1, $2, $3, $4) RETURNING id", request.Datetime, request.ParentCertificate.Id, request.Certificate.Id, request.Status).Scan(&id)
	if err != nil {
		return nil, err
	}
	request.Id = id
	return request, nil
}

func (repository *RequestRepositoryImpl) UpdateRequest(request *model.Request) error {
	_, err := repository.DB.Exec("UPDATE requests SET datetime = $1, status = $2 WHERE id = $3", request.Datetime, request.Status, request.Id)
	return err
}

func (repository *RequestRepositoryImpl) DeleteRequest(id int) error {
	_, err := repository.DB.Exec("DELETE FROM requests WHERE id = $1", id)
	return err
}
