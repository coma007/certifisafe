package model

type User struct {
	Id        int32
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	IsAdmin   bool
}
