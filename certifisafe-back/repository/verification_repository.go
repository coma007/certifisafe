package repository

import (
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"database/sql"
	"errors"
)

var (
	ErrNoVerificationWithEmail = errors.New("no request for given email")
	ErrNoVerificationWithCode  = errors.New("no request for given code")
)

type IVerificationRepository interface {
	GetVerification(id int32) (model.Verification, error)
	DeleteVerification(id int32) error
	CreateVerification(id int32, user model.Verification) (model.Verification, error)
	GetVerificationByCode(code string) (model.Verification, error)
	GetVerificationByEmail(email string) (*model.Verification, error)
	UpdateVerification(id int32, req model.Verification) (model.Verification, error)
}

type InMemoryVerificationRepository struct {
	Verifications []model.Verification
	DB            *sql.DB
}

func NewInMemoryVerificationRepository(db *sql.DB) *InMemoryVerificationRepository {
	var verifications = []model.Verification{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InMemoryVerificationRepository{
		Verifications: verifications,
		DB:            db,
	}
}

func (i *InMemoryVerificationRepository) GetVerification(id int32) (model.Verification, error) {
	stmt, err := i.DB.Prepare("SELECT * FROM verifications WHERE id=$1")

	utils.CheckError(err)

	var r model.Verification
	err = stmt.QueryRow(id).Scan(r.Id, r.Email, r.Code)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
			panic(err)
		}
		return model.Verification{}, err

	}
	return r, nil
}

func (i *InMemoryVerificationRepository) UpdateVerification(id int32, req model.Verification) (model.Verification, error) {
	stmt, err := i.DB.Prepare("UPDATE verifications" +
		" SET email=$1, code=$2" +
		" WHERE id=$4")

	utils.CheckError(err)

	_, err = stmt.Exec(req.Email, req.Code, id)

	utils.CheckError(err)
	return req, nil
}

func (i *InMemoryVerificationRepository) DeleteVerification(id int32) error {
	stmt, err := i.DB.Prepare("DELETE FROM verifications WHERE id=$1")
	utils.CheckError(err)

	_, err = stmt.Exec(id)
	return err
}

func (i *InMemoryVerificationRepository) CreateVerification(id int32, user model.Verification) (model.Verification, error) {
	stmt, err := i.DB.Prepare("INSERT INTO verifications(email, code)" +
		" VALUES($1, $2)")

	utils.CheckError(err)

	_, err = stmt.Exec(user.Email, user.Code)

	utils.CheckError(err)
	return user, nil
}

func (i *InMemoryVerificationRepository) GetVerificationByCode(code string) (model.Verification, error) {
	stmt, err := i.DB.Prepare("SELECT * FROM verifications WHERE code=$1")

	utils.CheckError(err)

	var r model.Verification
	err = stmt.QueryRow(code).Scan(&r.Id, &r.Email, &r.Code)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Verification{}, ErrNoVerificationWithCode
		}
		return model.Verification{}, err

	}
	return r, nil
}

func (i *InMemoryVerificationRepository) GetVerificationByEmail(email string) (*model.Verification, error) {
	stmt, err := i.DB.Prepare("SELECT id, email, code FROM verifications WHERE email=$1")

	utils.CheckError(err)

	var verification *model.Verification
	err = stmt.QueryRow(email).Scan(&verification.Id, &verification.Email, &verification.Code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoVerificationWithEmail
		}
		return nil, err
	}

	return verification, nil
}
