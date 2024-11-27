package model

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	FirstName *string
	LastName  *string
	Login     *string
	Phone     *string
	Email     *string
	Password  *string
	Role      *string
	IsOnline  *bool
	Status    *bool
}
