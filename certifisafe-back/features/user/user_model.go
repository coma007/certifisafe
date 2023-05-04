package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Deleted   gorm.DeletedAt
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	IsAdmin   bool
	IsActive  bool
}
