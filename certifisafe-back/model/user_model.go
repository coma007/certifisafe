package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id        int64 `gorm:"autoIncrement;PRIMARY_KEY"`
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	IsAdmin   bool
}
