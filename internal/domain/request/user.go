package request

type UpdateUser struct {
	FirstName *string `json:"first_name" form:"first_name"`
	LastName  *string `json:"last_name" form:"last_name"`
	Login     *string `json:"login" form:"login"`
	Phone     *string `json:"phone" form:"phone"`
	Email     *string `json:"email" form:"email"`
	Password  *string `json:"password " form:"password"`
}
