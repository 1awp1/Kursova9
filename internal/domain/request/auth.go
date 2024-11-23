package request

type Login struct {
	Login    string
	Password string
}

type Register struct {
	FirstName string
	LastName  string
	Login     string
	Phome     string
	Email     string
	Password  string
}
