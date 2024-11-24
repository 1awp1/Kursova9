package request

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Register struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     string `json:"login"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
