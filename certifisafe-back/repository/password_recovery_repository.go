package repository

import (
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"database/sql"
	"errors"
	"log"
)

var (
	ErrNoRequestWithEmail = errors.New("no request for given email")
	ErrNoRequestWithCode  = errors.New("no request for given code")
)

type IPasswordRecoveryRepository interface {
	GetRequest(id int32) (model.PasswordRecoveryRequest, error)
	DeleteRequest(id int32) error
	CreateRequest(id int32, user model.PasswordRecoveryRequest) (model.PasswordRecoveryRequest, error)
	GetRequestByCode(code string) (model.PasswordRecoveryRequest, error)
	GetRequestsByEmail(email string) ([]*model.PasswordRecoveryRequest, error)
}

type InMemoryPasswordRecoveryRepository struct {
	Requests []model.PasswordRecoveryRequest
	DB       *sql.DB
}

func NewInMemoryPasswordRecoveryRepository(db *sql.DB) *InMemoryPasswordRecoveryRepository {
	var requests = []model.PasswordRecoveryRequest{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InMemoryPasswordRecoveryRepository{
		Requests: requests,
		DB:       db,
	}
}

func (i *InMemoryPasswordRecoveryRepository) GetRequest(id int32) (model.PasswordRecoveryRequest, error) {
	stmt, err := i.DB.Prepare("SELECT * FROM passwordRecovery WHERE id=$1")

	utils.CheckError(err)

	var r model.PasswordRecoveryRequest
	err = stmt.QueryRow(id).Scan(r.Id, r.Email, r.Code)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
			panic(err)
		}
		return model.PasswordRecoveryRequest{}, err

	}
	return r, nil
}

func (i *InMemoryPasswordRecoveryRepository) DeleteRequest(id int32) error {
	stmt, err := i.DB.Prepare("DELETE FROM passwordRecovery WHERE id=$1")
	utils.CheckError(err)

	_, err = stmt.Exec(id)
	return err
}

func (i *InMemoryPasswordRecoveryRepository) CreateRequest(id int32, user model.PasswordRecoveryRequest) (model.PasswordRecoveryRequest, error) {
	stmt, err := i.DB.Prepare("INSERT INTO passwordRecovery(email, code)" +
		" VALUES($1, $2)")

	utils.CheckError(err)

	_, err = stmt.Exec(user.Email, user.Code)

	utils.CheckError(err)
	return user, nil
}

func (i *InMemoryPasswordRecoveryRepository) GetRequestByCode(code string) (model.PasswordRecoveryRequest, error) {
	stmt, err := i.DB.Prepare("SELECT * FROM passwordRecovery WHERE code=$1")

	utils.CheckError(err)

	var u model.PasswordRecoveryRequest
	err = stmt.QueryRow(code).Scan(&u.Id, &u.Email, &u.Code)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.PasswordRecoveryRequest{}, ErrNoRequestWithCode
		}
		return model.PasswordRecoveryRequest{}, err

	}
	return u, nil
}

func (i *InMemoryPasswordRecoveryRepository) GetRequestsByEmail(email string) ([]*model.PasswordRecoveryRequest, error) {
	stmt, err := i.DB.Prepare("SELECT id, email, code FROM passwordRecovery WHERE email=$1")

	utils.CheckError(err)

	var requests []*model.PasswordRecoveryRequest
	//err = stmt.QueryRow(email).Scan(&u.Id, &u.Email, &u.Code)
	queryLength := 0
	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r model.PasswordRecoveryRequest
		err := rows.Scan(&r.Id, &r.Email, &r.Code)
		if err != nil {
			log.Fatal(err)
		}
		requests = append(requests, &r)
		queryLength += 1

	}

	if queryLength == 0 {
		return nil, ErrNoRequestWithEmail
	}

	return requests, nil
}
