package request

type GetUsers struct {
	Page      *int
	Limit     *int
	Email     *string
	FirstName *string
	LastName  *string
	Role      *string
}
