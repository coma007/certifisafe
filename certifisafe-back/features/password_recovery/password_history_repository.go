package password_recovery

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrNoHistoryWithEmail = errors.New("no request for given email")
)

type PasswordHistoryRepository interface {
	GetHistory(id int32) (PasswordHistory, error)
	DeleteHistory(id int32) error
	CreateHistory(id int32, history PasswordHistory) (PasswordHistory, error)
	GetHistoryByEmail(email string) ([]*PasswordHistory, error)
	UpdateHistory(id int32, history PasswordHistory) (PasswordHistory, error)
}

type DefaultPasswordHistoryRepository struct {
	DB *gorm.DB
}

func NewDefaultPasswordHistoryRepository(db *gorm.DB) *DefaultPasswordHistoryRepository {
	return &DefaultPasswordHistoryRepository{
		DB: db,
	}
}

func (repository *DefaultPasswordHistoryRepository) CreateHistory(id int32, history PasswordHistory) (PasswordHistory, error) {
	result := repository.DB.Create(&history)
	return history, result.Error
}

func (repository *DefaultPasswordHistoryRepository) GetHistory(id int32) (PasswordHistory, error) {
	var r PasswordHistory
	result := repository.DB.First(&r, id)
	return r, result.Error
}

func (repository *DefaultPasswordHistoryRepository) GetHistoryByEmail(email string) ([]*PasswordHistory, error) {
	var history []*PasswordHistory
	result := repository.DB.Where("user_email=?", email).Find(&history)

	return history, result.Error
}

func (repository *DefaultPasswordHistoryRepository) UpdateHistory(id int32, req PasswordHistory) (PasswordHistory, error) {
	result := repository.DB.Save(&req)
	return req, result.Error
}

func (repository *DefaultPasswordHistoryRepository) DeleteHistory(id int32) error {
	result := repository.DB.Delete(&PasswordHistory{}, id)
	return result.Error
}
