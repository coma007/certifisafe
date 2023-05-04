package auth

import "gorm.io/gorm"

type Verification struct {
	gorm.Model
	Deleted gorm.DeletedAt
	Email   string
	Code    string
}
