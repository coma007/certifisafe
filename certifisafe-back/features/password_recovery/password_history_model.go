package password_recovery

import "gorm.io/gorm"

type PasswordHistory struct {
	gorm.Model
	Deleted           gorm.DeletedAt
	UserEmail         string
	ForbiddenPassword string
}
