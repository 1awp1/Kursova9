package inject

import "dim_kurs/internal/handler.go"

type Handlers struct {
	handler.IAuth
	handler.IUser
	handler.IMiddleware
}

func NewHandlers(usecases *UseCases) *Handlers {
	return &Handlers{
		IAuth:       handler.NewAuth(usecases.IAuth),
		IUser:       handler.NewUser(usecases.IUser),
		IMiddleware: handler.NewMiddleware(usecases.IAuth),
	}
}
