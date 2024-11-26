package request

type UpdateUser struct {
	FirstName *string `json:"first_name" form:"first_name" binding:"required"`
	LastName  *string `json:"last_name" form:"last_name" binding:"required"`
	Login     *string `json:"login" form:"login" binding:"required"`
	Phone     *string `json:"nhone" form:"phome" binding:"required"`
	Email     *string `json:"email" form:"email" binding:"required"`
	Password  *string `json:"password " form:"password " binding:"required"`
}
