package repository

import (
	"certifisafe-back/model"
	"certifisafe-back/utils"
	"database/sql"
	"errors"
)

var (
	ErrNoUserWithEmail = errors.New("no user for given email")
)

type IUserRepository interface {
	UpdateUser(id int32, user model.User) (model.User, error)
	GetUser(id int32) (model.User, error)
	DeleteUser(id int32) error
	CreateUser(id int32, user model.User) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
}

type InMemoryUserRepository struct {
	Users []model.User
	DB    *sql.DB
}

func NewInMemoryUserRepository(db *sql.DB) *InMemoryUserRepository {
	var users = []model.User{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return &InMemoryUserRepository{
		Users: users,
		DB:    db,
	}
}

func (i *InMemoryUserRepository) GetUser(id int32) (model.User, error) {
	stmt, err := i.DB.Prepare("SELECT * FROM users WHERE id=$1")

	utils.CheckError(err)

	var u model.User
	err = stmt.QueryRow(id).Scan(u.Id, u.Email, u.Password, u.FirstName, u.LastName, u.Phone, u.IsAdmin, u.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
			panic(err)
		}
		return model.User{}, err

	}
	return u, nil
}

func (i *InMemoryUserRepository) UpdateUser(id int32, user model.User) (model.User, error) {
	stmt, err := i.DB.Prepare("UPDATE users" +
		" SET email=$1, password=$2, first_name=$3, last_name=$4, phone=$5, is_admin=$6, is_active=$7" +
		" WHERE id=$8")

	utils.CheckError(err)

	_, err = stmt.Exec(user.Email, user.Password, user.FirstName, user.LastName, user.Phone, user.IsAdmin, user.IsActive, user.Id)

	utils.CheckError(err)
	return user, nil
}

func (i *InMemoryUserRepository) DeleteUser(id int32) error {
	stmt, err := i.DB.Prepare("DELETE FROM table_name WHERE id=$1")
	utils.CheckError(err)

	_, err = stmt.Exec(id)
	return err
}

func (i *InMemoryUserRepository) CreateUser(id int32, user model.User) (model.User, error) {
	stmt, err := i.DB.Prepare("INSERT INTO users(email, password, first_name, last_name, phone, is_admin, is_active)" +
		" VALUES($1, $2, $3, $4, $5, $6, $7)")

	utils.CheckError(err)

	_, err = stmt.Exec(user.Email, user.Password, user.FirstName, user.LastName, user.Phone, user.IsAdmin, user.IsActive)

	utils.CheckError(err)
	return user, nil
}

func (i *InMemoryUserRepository) GetUserByEmail(email string) (model.User, error) {
	stmt, err := i.DB.Prepare("SELECT * FROM users WHERE email=$1")

	utils.CheckError(err)

	var u model.User
	err = stmt.QueryRow(email).Scan(&u.Id, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Phone, &u.IsAdmin, &u.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrNoUserWithEmail
		}
		return model.User{}, err

	}
	return u, nil
}
