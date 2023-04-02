package repository

import (
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrRequestNotFound = errors.New("FromRepository - certificate not found")
)

type IRequestRepository interface {
	UpdateRequest(id int32, certificate model.Certificate) (model.Certificate, error)
	GetRequest(id int32) (model.Certificate, error)
	DeleteRequest(id int32) error
	CreateRequest(id int32, certificate model.Certificate) (model.Certificate, error)
}

type InmemoryRequestRepository struct {
	Requests []model.Certificate
	DB       *sql.DB
}

func NewInMemoryRequestRepository(db *sql.DB) *InmemoryRequestRepository {
	//var requests = new(model.Request);

	return &InmemoryRequestRepository{
		Requests: nil,
		DB:       db,
	}
}

func (i *InmemoryRequestRepository) UpdateReqyest(id int32, request model.Request) (model.Request, error) {
	for k := 0; k < len(i.Requests); k++ {
		if i.Requests[k].Id == id {
			return model.Request{}, nil
		}
	}

	return model.Request{}, nil
}

func (i *InmemoryRequestRepository) GetRequest(id int32) (model.Request, error) {
	rows, err := i.DB.Query("SELECT * FROM requests WHERE id=$1")
	utils.CheckError(err)

	var request model.Request

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		var parent_certificate_pk int
		var certificate_pk int
		var datetime time.Time
		var status model.RequestStatus
		err = rows.Scan(&id, &parent_certificate_pk, &certificate_pk, &datetime, &status)
		if err != nil {
			panic(err)
		}
		request := model.Request{id, nil, nil, datetime, status}
		fmt.Println(request)
		return request, err
	}
	return request, nil
}

func (i *InmemoryRequestRepository) DeleteRequest(id int32) error {
	for k := 0; k < len(i.Requests); k++ {
		if i.Requests[k].Id == id {
			// i.Certificates[k].Title = movie.Title
			return nil
		}
	}

	return nil
}

func (i *InmemoryRequestRepository) CreateRequest(id int32, request model.Request) (model.Request, error) {
	for k := 0; k < len(i.Requests); k++ {
		if i.Requests[k].Id == id {
			// i.Certificates[k].Title = movie.Title
			return model.Request{}, nil
		}
	}

	return model.Request{}, nil
}
