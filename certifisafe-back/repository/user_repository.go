package repository

import (
	"certifisafe-back/model"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNoUserWithEmail = errors.New("no user for given email")
)

type IUserRepository interface {
	UpdateUser(id int32, user model.User) (model.User, error)
	GetUser(id int32) (model.User, error)
	DeleteUser(id int32) error
	CreateUser(user model.User) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
}

type InMemoryUserRepository struct {
	Users []model.User
	DB    *gorm.DB
}

func NewInMemoryUserRepository(db *gorm.DB) *InMemoryUserRepository {
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
	var u model.User
	result := i.DB.First(&u, id)
	return u, result.Error
}

func (i *InMemoryUserRepository) UpdateUser(id int32, user model.User) (model.User, error) {
	result := i.DB.Save(&user)
	return user, result.Error
}

func (i *InMemoryUserRepository) DeleteUser(id int32) error {
	//TODO do this with deleted timestamp, see delete flag for gorm
	result := i.DB.Delete(&model.User{}, id)
	return result.Error
}

func (i *InMemoryUserRepository) CreateUser(user model.User) (model.User, error) {
	result := i.DB.Create(user)
	return user, result.Error
}

func (i *InMemoryUserRepository) GetUserByEmail(email string) (model.User, error) {
	var u model.User
	result := i.DB.Where("email=?", email).First(&u)
	return u, result.Error
}
