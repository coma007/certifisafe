package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Deleted         gorm.DeletedAt
	Email           string `gorm:"uniqueIndex"`
	Password        string
	LastPasswordSet time.Time
	FirstName       string
	LastName        string
	Phone           string
	IsAdmin         bool
	IsActive        bool
}
