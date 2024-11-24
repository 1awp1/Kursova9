package inject

import "dim_kurs/internal/handler.go"

type Handlers struct {
	handler.IAuth
}

func NewHandlers(usecases *UseCases) *Handlers {
	return &Handlers{
		IAuth: handler.NewAuth(usecases.IAuth),
	}
}
