package model

type User struct {
	Id        int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	IsAdmin   bool
}
