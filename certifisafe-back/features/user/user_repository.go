package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	UpdateUser(id int32, user User) (User, error)
	GetUser(id int32) (User, error)
	DeleteUser(id int32) error
	CreateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
}

type DefaultUserRepository struct {
	DB *gorm.DB
}

func NewDefaultUserRepository(db *gorm.DB) *DefaultUserRepository {
	return &DefaultUserRepository{
		DB: db,
	}
}

func (i *DefaultUserRepository) CreateUser(user User) (User, error) {
	result := i.DB.Create(&user)
	return user, result.Error
}

func (i *DefaultUserRepository) GetUser(id int32) (User, error) {
	var u User
	result := i.DB.First(&u, id)
	return u, result.Error
}

func (i *DefaultUserRepository) GetUserByEmail(email string) (User, error) {
	var u User
	result := i.DB.Where("email=?", email).First(&u)
	return u, result.Error
}

func (i *DefaultUserRepository) UpdateUser(id int32, user User) (User, error) {
	result := i.DB.Save(&user)
	return user, result.Error
}

func (i *DefaultUserRepository) DeleteUser(id int32) error {
	result := i.DB.Delete(&User{}, id)
	return result.Error
}
