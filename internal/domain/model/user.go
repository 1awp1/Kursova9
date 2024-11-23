package model

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Login     string
	Phone     string
	Email     string
	Password  string
	Role      string
	Token     string
	Status    bool
}
