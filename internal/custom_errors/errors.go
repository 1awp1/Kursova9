package custom_errors

import "errors"

var (
	AlreadyExist = errors.New("the user is already registered")
	UserNotExist = errors.New("the user does not exist yet")
	RoleNotExist = errors.New("the role does not exist yet")
)
