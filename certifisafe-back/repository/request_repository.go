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
	CreateRequest(id int32, request model.Request) (model.Request, error)
	UpdateRequest(id int32, request model.Request) (model.Request, error)
	DeleteRequest(id int32) error
}

type RequestRepositoryImpl struct {
	DB *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepositoryImpl {
	return &RequestRepositoryImpl{
		DB: db,
	}
}

func (repository *RequestRepositoryImpl) GetRequest(id int) (*model.Request, error) {
	request := &model.Request{}
	// TODO add parent certificate and certificate
	err := repository.DB.QueryRow("SELECT id, datetime, status FROM requests WHERE id = $1", id).Scan(&request.Id, &request.Datetime, &request.Status)

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
